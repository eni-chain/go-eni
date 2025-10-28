// SPDX-License-Identifier: GPL-3.0

pragma solidity >= 0.8.0;

import "./common.sol";
import "./localLog.sol";
import "./delegateCallBase.sol";


contract StakeManager is DelegateCallBase, Common {
    //锁仓周期
    uint internal constant LOCK_TIME_90 = 90 days;
    uint internal constant LOCK_TIME_180 = 180 days;
    uint internal constant LOCK_TIME_360 = 360 days;
    uint internal constant LOCK_TIME_720 = 720 days;
    uint internal constant LOCK_TIME_1440 = 1440 days;
    //uint internal constant MIN_PLEDGE_AMOUNT = 10000000000000000000000; //wei
    //uint internal constant MAX_PLEDGE_AMOUNT = 1000000000000000000000000; //wei

    //质押算力根据锁仓周期放大的倍数
    uint internal constant POW_MULTI_ERR = 0;
    uint internal constant POW_MULTI1 = 1;
    uint internal constant POW_MULTI2 = 2;
    uint internal constant POW_MULTI3 = 3;
    uint internal constant POW_MULTI4 = 4;
    uint internal constant POW_MULTI5 = 5;

    //续质押参数，1为不再续质押，2为只续质押本金，3为之前的本金和利息一起续质押
    uint internal constant NO_CONTINUE = 1;
    uint internal constant CONTINUE_PLEDGE = 2;
    uint internal constant COMPOUND_INTEREST = 3;

    //冷却时间，避免用户频繁执行高耗能的操作，比如投票者频繁切换验证者，增加链上负担
    uint256 internal coolingPeriod = 7 days;

    //退出缓冲期，用以检查验证者在退出前是否作恶，作恶则罚没质押金
    uint256 internal exitingPeriod = 14 days;

    //每个区块的总质押奖励，单位为wei，总量为0.5ENI
    uint256 internal totalRewardPerBlock = 600000000000000000;

    //以下三个状态变量是用于流动性挖矿计算的全局参数
    uint256 internal rewardPerShare; //每股奖励
    uint256 internal lastRewardBlock;//上次奖励区块
    uint256 internal totalStaked;    //全部奖励质押份数，这是根据质押周期调整算力倍数后的值

    struct OrderId { //可以试着将持有者和编号用abi.encodePacked打包成一个byte32替换
        address shareholder;   //质押订单持有者
        uint256 sequence;      //质押订单编号
    }

    struct Order {
      uint256 amount;               //质押额
      uint256 enterTime;            //质押开始时间
      uint256 lockPeriod;           //锁仓周期，根据锁仓周期和质押额，可计算质押倍率后的算力
      uint256 rewardDebt;           //奖励债务，用于记录已结算过的奖励
      uint256 unclaimedReward;      //未提取奖励，用于记录已结算未提取的奖励
      uint256 continueFlag;         //续质押标识，锁仓到期继续质押
      uint256 coolingExpired;       //切换验证者冷却到期时间
      address validator;            //指向的验证者，未指向为address(0)
      //int256 currentValidatorIdx;   //当前指向验证者的索引，如果为-1，表示未指向验证者
    }

    struct Voter {
        Order[] orders; //一个投票者可以持有多笔质押订单,通过编号查询具体订单
    }

    struct Validator {
        Order order;        //验证者的质押订单,其指向的验证者为空地址
        uint256 poll;       //总得票数，验证者自质押额+总得票额可计算出块奖励倍数
        OrderId[] voters;   //投票者列表, 记录哪些质押指向了当前验证者。验证者退出时，用户找到投票者，触发其切换指向的验证者。
        uint256 exitExpired;//退出缓冲期截止日期
        bool frozen;        //验证者质押是否已冻结/罚没
    }

    //单个锁仓周期的订单时序表，相同锁仓周期的订单放在一个表内，订单的插入时间和到期时间都是顺序同频的
    struct OrderSeq {
        uint256 realStartIdx;//订单都是按顺序从头部开始删除，虽然数据从DB删除了，但是槽位还在，遍历时需要跳过，因此要记录第一个有实际值的下标
        OrderId[] seqList;
    }

    //验证者转手信息
    struct validatorChange {
        address receiver;
        address node;
        bytes pubKey;
        uint256 expiredTime;
    }

    //质押订单表内有五个分表，每种锁仓周期，各有一个单独订单时序表
    OrderSeq[] internal _orderSequences;

    //投票者质押表, 用OrderId.shareholder检索投票者，用OrderId.sequence查找投票者的订单
    mapping(address => Voter) internal _voters;

    //验证者质押表, 用OrderId.shareholder检索验证者
    mapping (address => Validator) internal _validators;

    //待退出验证者索引表
    address[] _exitings;

    //待退出验证者列表
    mapping (address=>uint256) internal _exitingValidators;

    //待转手验证者索引表
    address[] _changes;

    //验证者转手列表
    mapping(address=>validatorChange) internal _changeValidators;

    //更新奖励池
    function updatePool() internal {
        // 1. 如果当前区块还没到上次更新的区块，直接返回
        if (block.number <= lastRewardBlock){
            return;
        }

        // 2. 如果总质押量为0（还没人质押），则只更新区块号，不发放奖励
        if (totalStaked == 0) {
            lastRewardBlock = block.number;
            return;
        }

        // 3. 计算经过了几个区块
        uint256 blocksPassed = block.number - lastRewardBlock;

        // 4. 计算这段时间内产生的总奖
        uint256 totalReward = blocksPassed * totalRewardPerBlock;

        // 5. 计算每份额应增加的奖励，一个eni为1份(不是1wei)
        uint256 rewardPerSharePhased = totalReward / (totalStaked / 1e18);

        // 6. 更新全局累计每份额奖励
        rewardPerShare += rewardPerSharePhased;

        // 7. 更新最后一次奖励区块为当前区块
        lastRewardBlock = block.number;
    }

    function claimRewardBasic(Order storage order) internal {
        //应该用order.amount先除1e18得出有多少股份，然后再用股份数乘rewardPerShare
        order.unclaimedReward += ((order.amount / 1e18) * rewardPerShare) - order.rewardDebt;
        order.rewardDebt = (order.amount / 1e18) * rewardPerShare;
    }

    //投票者领取奖励
    function claimReward(uint256 sequence) external {
        Voter storage voter = _voters[msg.sender];
        require(voter.orders.length > 0, "There is no order for current user");
        require(sequence < voter.orders.length, "Invalid order sequence");

        Order storage order = voter.orders[sequence];
        require(order.amount != 0, "There is no order for the sequence");

        //require(block.timestamp <= order.coolingExpired, "The cooling-off period has not expired");
        //order.enterTime+order.lockPeriod>=block.timestamp时，会自动计算奖励并根据续质押标识处理复投，并重置enterTime，所以，超期一定未续质押，不再计算奖励
        require(order.enterTime + order.lockPeriod < block.timestamp, "The lock-up period expires but the pledge is not renewed");

        //未指定验证者，不计算奖励
        require(order.validator != address(0), "There is no validator for current order");
        claimRewardBasic(order);
    }

    //验证者领取奖励：因为验证者只有一个质押订单，无需序号
    function claimReward() external {
        Validator storage val = _validators[msg.sender];
        require(val.order.amount != 0, "msg.sender is not validator");
        require(val.frozen == false, "validator was frozen");

        Order storage order = val.order;
        require(order.enterTime + order.lockPeriod < block.timestamp, "The lock-up period expires but the pledge is not renewed");

        claimRewardBasic(order);
    }

    //根据锁仓周期推到质押放大倍数
    function multiNumber(uint256 lockPeriod) internal pure returns(uint256){
        if(lockPeriod == LOCK_TIME_90){
            return POW_MULTI1;
        }else if(lockPeriod == LOCK_TIME_180){
            return POW_MULTI2;
        }else if(lockPeriod == LOCK_TIME_360){
            return POW_MULTI3;
        }else if(lockPeriod == LOCK_TIME_720){
            return POW_MULTI4;
        }else if(lockPeriod == LOCK_TIME_1440){
            return POW_MULTI5;
        }else{
            return POW_MULTI_ERR;
        }
    }

    //续质押合法值检查
    function continueFlagValid(uint256 continueFlag) internal pure returns(bool){
        if((continueFlag == NO_CONTINUE) ||(continueFlag == CONTINUE_PLEDGE)|| (continueFlag == COMPOUND_INTEREST)){
            return true;
        }
        return false;
    }

    //填充基本订单
    function fillOrder(Order storage order, uint256 lockPeriod, uint256 continueFlag) internal{
        order.amount = msg.value;
        order.enterTime = block.timestamp;
        order.lockPeriod = lockPeriod;
        order.rewardDebt = (order.amount / 1e18) * rewardPerShare;
        order.unclaimedReward = 0;
        order.continueFlag = continueFlag;
        order.coolingExpired = block.timestamp + coolingPeriod;
    }

    //验证者质押
    function validatorStake(uint256 lockPeriod, uint256 continueFlag, address node, bytes calldata pubKey) payable external {
        uint256 multi = multiNumber(lockPeriod);
        require(multi != POW_MULTI_ERR, "Lock period error!");
        require(msg.value >= MIN_PLEDGE_AMOUNT, "The transfer amount is less than validator minimum pledge amount!");
        require(continueFlagValid(continueFlag), "Continue flag error!");

        updatePool();
        Validator storage vali = _validators[msg.sender];
        require(vali.order.amount == 0, "Validator alread exist.");

        Order storage order = vali.order;
        fillOrder(order, lockPeriod, continueFlag);

        uint256 pledgePower = order.amount * multi;
        totalStaked += pledgePower; //将放大后的算力加总到总股本中

        IValidatorManager(VALIDATOR_MANAGER_ADDR).addValidator(
            msg.sender,
            node,
            msg.sender,
            msg.value,
            block.number,
            "",
            "",
            pubKey
        );

        //将订单ID插入时序表中，供自动化检测到期处理
        OrderId memory orderId;
        orderId.shareholder = msg.sender;

        OrderSeq storage orderSeq = _orderSequences[multi];
        orderSeq.seqList.push(orderId);
    }

    //投票者质押
    function voteStake(uint256 lockPeriod, uint256 continueFlag, address validator) payable external{
        uint256 multi = multiNumber(lockPeriod);
        require(multi != POW_MULTI_ERR, "Lock period error!");
        require(msg.value >= 1e18, "The transfer amount is less than 1 ENI!");
        require(validator != address(0), "Validator address is null!");

        //检查指向的验证者是为合法验证者
        Validator storage vali = _validators[validator];
        require(vali.order.amount != 0, "Validator not exist.");

        updatePool();

        Voter storage voter = _voters[msg.sender];
        Order storage order = voter.orders[voter.orders.length];

        fillOrder(order, lockPeriod, continueFlag);
        order.validator = validator;

        //根据锁仓周期放大质押算力
        uint256 pledgePower = order.amount * multi;
        totalStaked += pledgePower;

        //验证者的投票额加上投票者的质押额
        vali.poll += msg.value;

        //验证者中记录支持自己的订单
        OrderId memory orderId;
        orderId.shareholder = msg.sender;
        orderId.sequence = voter.orders.length;
        vali.voters.push(orderId);

        //将订单ID插入时序表中，供自动化检测到期处理
        OrderSeq storage orderSeq = _orderSequences[multi];
        orderSeq.seqList.push(orderId);
    }

    function validatorValid(Validator storage vali) internal view returns (bool) {
        if(vali.frozen == false && vali.exitExpired == 0 && vali.order.amount != 0){
            return true;
        }
        return false;
    }

    function markValidatorExit(address validator) internal {
        //1.调用验证者管理合约为节点打上退出标签
        IValidatorManager(VALIDATOR_MANAGER_ADDR).markExitValidator(validator);

        //2.将验证者插入待退出验证者列表（key:验证者地址，value:缓冲截止日期）
        _exitings.push(validator);
        _exitingValidators[validator] = block.timestamp + exitingPeriod;
    }

    function markValidatorTransfer(address from, address to, address node, bytes memory pk) internal {
        //1.调用验证者管理合约为节点打上退出标签
        IValidatorManager(VALIDATOR_MANAGER_ADDR).markExitValidator(from);

        //2.生成验证者转手信息，将验证者转手信息以转出方地址为key存储
        _changes.push(from);
        validatorChange storage ch = _changeValidators[from];
        ch.receiver = to;
        ch.node = node;
        ch.pubKey = pk;
        ch.expiredTime = block.timestamp + exitingPeriod;
    }

    //缓冲期到期后，被区块自动调用
    function validatorExit(address validator, uint256 amount) internal {
        //1.调用验证者管理合约删除验证者信息
        // --验证者管理合约会在下轮epoch自动删除有退出标记的验证者，所以这里不用调用了

        //2.遍历投票者信息，更新投票者奖励，将投票者指向的验证者删除
        Validator storage v = _validators[validator];
        for(uint i = 0; i < v.voters.length; i++){
            address voterAddr = v.voters[i].shareholder;
            uint seq = v.voters[i].sequence;

           Voter storage user = _voters[voterAddr];
           Order storage order = user.orders[seq];
           order.validator = address(0);
        }

        //3.返还质押金和奖励
        payable(validator).transfer(amount);

        //4.删除验证者相关信息和质押订单
        for(uint i = 0; i < _exitings.length; i++){
            if(_exitings[i] == validator){
                _exitings[i] = _exitings[_exitings.length-1];
                _exitings.pop();
            }
        }
        delete _exitingValidators[validator];
        delete _validators[validator];

    }

    //缓冲期到期后，被区块自动调用
    function validatorTransfer(address from) internal {
        //1.调用验证者管理合约删除转出验证者信息
        // --验证者管理合约会在下轮epoch自动删除有退出标记的验证者，所以这里不用调用了

        //2.调用验证者管理合约为接收者生成验证者信息
        Validator storage v = _validators[from];
        validatorChange storage ch = _changeValidators[from];
        IValidatorManager(VALIDATOR_MANAGER_ADDR).addValidator(
            ch.receiver,
            ch.node,
            ch.receiver,
            v.order.amount,
            block.number,
            "",
            "",
            ch.pubKey
        );

        //3.遍历投票者信息，将投票者指向的验证者更新
        for(uint i = 0; i < v.voters.length; i++){
            address voterAddr = v.voters[i].shareholder;
            uint seq = v.voters[i].sequence;

           Voter storage user = _voters[voterAddr];
           Order storage order = user.orders[seq];
           order.validator = from;
        }

        //4.将原验证者的质押订单准到新验证者地址下，并更新查询关系
        _validators[ch.receiver] = v;

        //5.返还转出验证者的质押金和奖励
        //转手交易只负责转手，转手的资金由用户线下处理？
         //payable(from).transfer(v.order.amount);

        //6.删除验证者信息和质押订单
        for(uint i = 0; i < _changes.length; i++){
            if(_changes[i] == from){
                _changes[i] = _changes[_changes.length-1];
                _changes.pop();
            }
        }
        delete _changeValidators[from];
        delete _validators[from];
    }



    // 自动处理内容(首先更新奖励池计算奖励):
    // - 遍历验证者待退出列表，缓冲期到期的：将验证者删除，并返还质押，同时遍历所有投票者，更新其奖励，并将其指向验证者删除。
    function autoProcByBlock() external {

        //订单到期处理
        for(uint i = 0; i < _orderSequences.length; i++){
            //对每个锁仓周期的订单时序表分别处理
            OrderSeq storage orderSeq = _orderSequences[i];

            for(uint ii = orderSeq.realStartIdx; ii < orderSeq.seqList.length; ii++){
                Order storage order;
                OrderId storage id = orderSeq.seqList[ii];
                if(_voters[id.shareholder].orders.length > 0){
                    //投票者持有的质押订单
                    require(id.sequence < _voters[id.shareholder].orders.length, "invalid order sequence");
                    order = _voters[id.shareholder].orders[id.sequence];

                    if(order.continueFlag == CONTINUE_PLEDGE){
                        //复投
                        order.enterTime = block.timestamp;
                    }else if(order.continueFlag == COMPOUND_INTEREST){
                        //复利
                        Validator storage vali = _validators[order.validator];
                        if(validatorValid(vali)){
                            //验证者还在，且未被冻结
                            if(vali.order.amount + vali.poll + order.unclaimedReward > MAX_PLEDGE_AMOUNT){
                                //复利超出验证者总算力限额
                                //计算最终奖励
                                claimRewardBasic(order);

                                //复利导致验证者总质押超额，将质押额从验证者得票额中减除
                                vali.poll -= order.amount;

                                //将订单指向的验证者索引设为无效值
                                order.validator = address(0);

                                //从订单有序表中删除订单ID
                                delete orderSeq.seqList[ii];
                                orderSeq.realStartIdx += 1;

                                //将投票订单从验证者的投票者列表中删除
                                for(uint iii = 0; iii<vali.voters.length; iii++){
                                    if(vali.voters[iii].shareholder == id.shareholder && vali.voters[iii].sequence == id.sequence){
                                        delete vali.voters[iii];
                                    }
                                }
                            }else {
                                //复利未导致验证者总质押额超额
                                order.amount += order.unclaimedReward;
                                order.enterTime = block.timestamp;
                                vali.poll += order.unclaimedReward;
                            }
                        }
                    }else{
                        //不续，到期赎回
                        //计算最终奖励
                        claimRewardBasic(order);

                        //更新验证者的投票列表和得票额
                        Validator storage vali = _validators[order.validator];
                        if(validatorValid(vali)){
                            //验证者得票额减去订单质押额
                            vali.poll -= order.amount;
                            //将订单从验证者列表中删除
                            for(uint iii = 0; iii<vali.voters.length; iii++){
                                if(vali.voters[iii].shareholder == id.shareholder && vali.voters[iii].sequence == id.sequence){
                                    delete vali.voters[iii];
                                }
                            }
                        }

                        //从订单有序表中删除订单ID
                        delete orderSeq.seqList[ii];
                        orderSeq.realStartIdx += 1;

                        //返还质押金额
                        payable(id.shareholder).transfer(order.amount+order.unclaimedReward);
                    }
                }else{
                    //验证者持有的质押订单
                    order = _validators[id.shareholder].order;
                    if(order.continueFlag == CONTINUE_PLEDGE){
                        //复投
                        order.enterTime = block.timestamp;
                    }else if(order.continueFlag == COMPOUND_INTEREST){
                        //复利
                        //todo: 验证者复利超过了质押额上限该如何处理？是否让验证者退出？还是按复投处理？此处暂按复投处理，后期需要讨论确认
                         Validator storage vali = _validators[id.shareholder];
                        if(order.amount + vali.poll + order.unclaimedReward < MAX_PLEDGE_AMOUNT){
                            order.amount += order.unclaimedReward;
                            order.enterTime = block.timestamp;
                        }else{
                            //复利超出质押额上限，按照复投处理
                            order.enterTime = block.timestamp;
                        }
                    }else{
                        //不续, 将验证者标记为待退出验证者
                        //计算最终奖励
                        claimRewardBasic(order);
                        markValidatorExit(id.shareholder);
                    }
                }
            }
        }

        //遍历待退出验证者列表，处理到期退出者
        for(uint i = 0; i < _exitings.length; i++){
            uint expired = _exitingValidators[_exitings[i]];
            if(block.timestamp >= expired){
                Validator storage vali = _validators[_exitings[i]];
                validatorExit(_exitings[i], vali.order.amount + vali.order.unclaimedReward);
            }
        }

        //遍历转手验证者列表，处理到期转手者
        for(uint i = 0; i < _changes.length; i++){
            validatorChange storage vc = _changeValidators[_changes[i]];
            if(block.timestamp >= vc.expiredTime){
                validatorTransfer(_changes[i]);
            }
        }
    }

    //转移质押
    function transferStake() external{

    }

    //续质押，修改续质押参数，选项有不续，复投，复利，其他为非法参数，直接报错返回
    function reStake(uint256 continueFlag) external{

    }

    //切换指定订单指向的验证者
    function changeOrderValidator(uint256 sequence, address validator) external{

    }

    //提取奖励
    function withdrawProfit() external{

    }

    //到期赎回
    function redemption() external{

    }

}