#!/bin/bash

# 检查文件参数是否存在
if [ $# -lt 1 ]; then
    echo "用法: $0 <交易文件.txt> [rpc-url]"
    echo "默认 rpc-url: http://localhost:8545"
    exit 1
fi

TX_FILE="$1"
RPC_URL="${2:-http://localhost:8545}"

# 验证文件是否存在
if [ ! -f "$TX_FILE" ]; then
    echo "错误：文件 $TX_FILE 未找到"
    exit 1
fi

# 设置失败计数器
fail_count=0
success_count=0
total=$(wc -l < "$TX_FILE")

# 逐行处理交易
while IFS= read -r raw_tx; do
    # 清理交易数据
    tx=$(echo "$raw_tx" | tr -d '[:space:]')

    # 基础格式验证
    if [[ ! $tx =~ ^0x[0-9a-fA-F]+$ ]]; then
        echo "跳过无效交易: ${tx:0:20}..."
        ((fail_count++))
        continue
    fi

    # 发送原始交易
    echo "正在发送交易 (${success_count}/${total}): ${tx:0:20}..."

    if output=$(cast rpc eth_sendRawTransaction "$tx" --rpc-url "$RPC_URL" 2>&1); then
        echo "成功 | 哈希: $output"
        ((success_count++))
    else
        echo "失败 | 错误: $output"
        ((fail_count++))
    fi

done < "$TX_FILE"

# 输出统计结果
echo "发送完成"
echo "成功: $success_count"
echo "失败: $fail_count"
echo "总计: $total"