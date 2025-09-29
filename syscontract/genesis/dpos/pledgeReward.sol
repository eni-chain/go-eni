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
    uint internal constant MAX_PLEDGE_AMOUNT = 1000000000000000000000000; //wei

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

    //切换验证者冷却时间，避免投票者频繁切换验证者，增加链上负担
    uint256 internal coolingPeriod = 7 days;

    //每个区块的总质押奖励
    uint256 internal totalRewardPerBlock = 600000000000000000; //单位为wei

    //以下三个状态变量是用于流动性挖矿计算的全局参数
    uint256 internal rewardPerShare; //每股奖励
    uint256 internal lastRewardBlock;//上次奖励区块
    uint256 internal totalStaked;    //全部奖励质押，这是根据质押周期调整算力倍数后的值

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
      int256 currentValidatorIdx;   //当前指向验证者的索引，如果为-1，表示未指向验证者
    }

    struct Voter {
        Order[] orders; //一个投票者可以持有多笔质押订单,通过编号查询具体订单
    }

    struct Validator {
        Order order;        //验证者的质押订单,其指向的验证者列表为空
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

    //质押订单表内有五个分表，每种锁仓周期，各有一个单独订单时序表
    OrderSeq[] internal OrderSequences_;

    //投票者质押表, 用OrderId.shareholder检索投票者，用OrderId.sequence查找投票者的订单
    mapping(address => Voter) internal Voters_;

    //验证者质押表, 用OrderId.shareholder检索验证者
    mapping (address => Validator) internal Validators_;

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
        //应该用order.amount先除1e18得出有多少股份，然后再用股份数乘rewardPerShare，但为了计算精度，采取了以下写法
        order.unclaimedReward += (order.amount * rewardPerShare / 1e18) - order.rewardDebt;
        order.rewardDebt = order.amount * rewardPerShare / 1e18;
    }

    //投票者领取奖励：因为验证者只有一个质押订单，无需序号
    function claimReward(uint256 sequence) external {
        Voter storage v = Voters_[msg.sender];
        require(v.orders.length > 0, "There is no order for current user");
        require(sequence < v.orders.length, "Invalid order sequence");

        Order storage order = v.orders[sequence];
        require(order.amount != 0, "There is no order for the sequence");

        //require(block.timestamp <= order.coolingExpired, "The cooling-off period has not expired");
        //order.enterTime+order.lockPeriod>=block.timestamp时，会自动计算奖励并根据续质押标识处理复投，并重置enterTime，所以，超期一定未续质押，不再计算奖励
        require(order.enterTime + order.lockPeriod < block.timestamp, "The lock-up period expires but the pledge is not renewed");

        //未指定验证者，不计算奖励
        require(order.validator != address(0), "There is no validator for current order");
        claimRewardBasic(order);
    }

    //验证者领取奖励
    function claimReward() external {
        Validator storage val = Validators_[msg.sender];
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
        order.rewardDebt = order.amount * rewardPerShare / 1e18;
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
        Validator storage vali = Validators_[msg.sender];
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

        OrderSeq storage orderSeq = OrderSequences_[multi];
        orderSeq.seqList.push(orderId);
    }

    //投票者质押
    function voteStake(uint256 lockPeriod, uint256 continueFlag, address validator) payable external{
        uint256 multi = multiNumber(lockPeriod);
        require(multi != POW_MULTI_ERR, "Lock period error!");
        require(msg.value >= 1e18, "The transfer amount is less than 1 ENI!");
        require(validator != address(0), "Validator address is null!");

        //检查指向的验证者是为合法验证者
        Validator storage vali = Validators_[validator];
        require(vali.order.amount != 0, "Validator not exist.");

        updatePool();

        Voter storage voter = Voters_[msg.sender];
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
        OrderSeq storage orderSeq = OrderSequences_[multi];
        orderSeq.seqList.push(orderId);
    }

    function validatorValid(Validator storage vali) internal view returns (bool) {
        if(vali.frozen == false && vali.exitExpired == 0 && vali.order.amount != 0){
            return true;
        }
        return false;
    }

    function autoProcByBlock() external {
        for(uint i = 0; i < OrderSequences_.length; i++){ //对每个锁仓周期的订单时序表分别处理
            OrderSeq storage orderSeq = OrderSequences_[i];

            for(uint ii = orderSeq.realStartIdx; ii < orderSeq.seqList.length; ii++){
                OrderId storage id = orderSeq.seqList[ii];
                bool isValidator;
                Order storage order;
                if(Voters_[id.shareholder].orders.length > 0){//为投票者持有的质押订单
                    isValidator = false;
                    require(id.sequence < Voters_[id.shareholder].orders.length, "invalid order sequence");
                    order = Voters_[id.shareholder].orders[id.sequence];
                }else{
                    isValidator = true;
                    //为验证者持有的质押订单
                    order = Validators_[id.shareholder].order;
                }

                if(block.timestamp >= (order.enterTime + order.lockPeriod)){
                    claimRewardBasic(order);
                    if(order.continueFlag == NO_CONTINUE){
                        //不续投
                        delete orderSeq.seqList[ii];//存储槽被保留，但DB不再存储数据，gas费将被退还
                        orderSeq.realStartIdx += 1; //如果realStartIdx==seqList.length,for循环不会进入处理，push新元素后，realStartIdx正好指向该元素位置，而length会+1
                    }else if(order.continueFlag == CONTINUE_PLEDGE){
                        //复投
                        order.enterTime = block.timestamp;
                    }else if(order.continueFlag == COMPOUND_INTEREST){
                        //复利
                        if(isValidator){
                            Validator storage vali = Validators_[id.shareholder];
                            if(validatorValid(vali)){
                                if(vali.order.amount + vali.order.unclaimedReward + vali.poll > MAX_PLEDGE_AMOUNT){
                                    //从订单时序表中删除订单ID
                                    delete orderSeq.seqList[ii];
                                    orderSeq.realStartIdx += 1; //如果realStartIdx==seqList.length,for循环不会进入处理，push新元素后，realStartIdx正好指向该元素位置，而length会+1
                                }else{
                                    order.enterTime = block.timestamp;
                                    order.amount += order.unclaimedReward;
                                }
                            }else{
                                //todo: 发送事件提醒验证者
                            }
                        }else{
                            Validator storage vali = Validators_[order.validator];
                            if(validatorValid(vali)){//验证者还在，且未被冻结
                                //复利超出验证者总算力限额
                                if(vali.order.amount + vali.poll + order.amount > MAX_PLEDGE_AMOUNT){
                                    //复利失败，将质押额从验证者得票额中减除
                                    vali.poll -= order.amount;
                                    for(uint iii = 0; iii<vali.voters.length; iii++){
                                        if(vali.voters[iii].shareholder == id.shareholder && vali.voters[iii].sequence == id.sequence){
                                            delete vali.voters[iii];
                                        }
                                    }
                                    //将订单指向的验证者索引设为无效值
                                    order.validator = address(0);

                                    //从订单有序表中删除订单ID
                                    delete orderSeq.seqList[ii];
                                    orderSeq.realStartIdx += 1;
                                }else {
                                    order.amount += order.unclaimedReward;
                                    order.enterTime = block.timestamp;
                                    vali.poll += order.unclaimedReward;
                                }
                            }else{

                                //todo: 发送事件提醒投票者
                            }
                        }
                    }else{
                        //非法参数
                        revert("invalid continue flag");
                    }
                }
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