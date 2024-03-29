# go_nexa

测试准备
* 将 `accounts_template.json` 文件重命名为 `accounts.json`
* 根据 `accounts.json`中格式配置自己的账户信息


官方文档
https://spec.nexa.org/protocol/blockchain/transaction/transaction-signing
https://spec.nexa.org/protocol/blockchain/transaction



## Previous Outputs Hash 

官方文档的描述是针对`bitcoincash`的描述

```
input_type = 00 必须是0
pre_out_raw = input_type1 + reverse(pre_outpoint_hex1) + input_type2 + reverse(pre_outpoint_hex2) + ........
previous_outputs_hash  = inputsHash = doubleSha256(pre_out_raw)
```

## inputsAmountHash
```
LE = littleEnd
input_amount_raw = LE(amount1) + LE(amount2) + .....
inputsAmountHash = doubleSha256(input_amount_raw)
```
## inputsSequenceHash
```
input_sequence_raw = LE(sequence1) + LE(sequence2) + ....
inputsAmountHash = doubleSha256(input_sequence_raw)
```

## script
```
script = LE(len(script)) + script_pub_key
当signtype == 0，即ALL时 script = 026cad

```

## outpusHash
```
output_type = 01 表示本币
ouput_raw = output_type1 + LE(out_amount1) + script_pub_key1 
+ output_type2 + LE(out_amount2) + script_pub_key2 
+ ....
outpusHash = doubleSha256(ouput_raw)
```



## Tx hash raw 格式
要按照如下顺序拼接，这和官方文档还是有出入的
```
00 version
b4b7bba1231afeee623d1b4c4a5bc556d190a9a1560b0cce4155389ccb382b38 inputsHash
24dc6b0060bf62a01f9fa31f3c81786b2329183d5c9338f6f91c82272b72eea6 inputsAmountHash
c992651ac89a97aecd0811c1761915a8e2c8f5153d1bdd994a789e6bd86ab717 inputsSequenceHash
026cad script
5bcc3e2441b73b3964434ca242be2ad446a84bba6bf08c6f2322fa592ca19629 outpusHash
56df0300 locktime
00 signtype
```


## 多个utxo的签名交易格式
```
00 version
02 input_count
00 outpoint_type1 必须是0
94a577802f958a3c151411f91ce434897b705964ed5ae56bbf26df0d314bc5b3 reverse(outpoint_1)
64 outpoint_sign_all_1
2221 pubkey_all_len_1
02c732230b0ae3cd0142508e3388e9eff47d063d3046ab5c9147d8e76b8bb03b71 pubkey_1
40 sing_len1
5dc2a01ff6963d50ef629548f557d282a6e1024607640690bfcd6b9fdc20321b sign_r_1
bec60d41b64ab42d9b1705112c2010aa9805ef1394f8e09ec77d0dac5c158c98 sign_s_1
feffffff sequence_1
1027000000000000  0x2710 = 10000 outpoint_amonut_1
00 00 outpoint_type2 必须是0
e431c41ddf202d5bd73c4cabf260c23f12cc4368d12630fff71a2badb5adf684 reverse(outpoint_2)
64 outpoint_sign_all_2
2221 pubkey_all_len_2
02c732230b0ae3cd0142508e3388e9eff47d063d3046ab5c9147d8e76b8bb03b71 pubkey_2
40 sing_len2
5dc2a01ff6963d50ef629548f557d282a6e1024607640690bfcd6b9fdc20321b sign_r_2
bec60d41b64ab42d9b1705112c2010aa9805ef1394f8e09ec77d0dac5c158c98 sign_s_2
feffffff sequence_2
c832000000000000  0x32c8 = 130000  outpoint_amonut_2
02 output_count
01 output_type
204e000000000000 output_amount1
170051141476bc2644608f984125fd4b7875f067bc5953b8 output_script_pubkey1
01 output_type
7107000000000000 output_amount2
170051141129a5ec6501c423c686247dfe7f413b4ebf7449 output_script_pubkey2
56df0300 locktime
```

