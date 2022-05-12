# **Overview**

这是一个使用 musig2的 api 文档。 这些有助于使用Go来构建阈值签名钱包。 为了应对taproot升级，这个api还提供了taproot普通交易和门槛钱包交易的构建。

# **Api**

## **Construct Transaction**

下面是构造交易相关的函数

---

### GenerateRawTx(prev_txs, txids, input_indexs, addresses, amounts)**

#### 说明

构建一个原始的交易，用于下面计算交易哈希然后签名。输入的交易原文，交易id和输入的交易索引必须一一对应。输出的地址和输出的数量必须一一对应。支持**op_return**,只需将amout设置为0，相应的address设置需要附带的信息即可。

#### **参数和返回值**

| **Name**         | **Type** | **Description**    |
| ---------------- | -------- | ------------------ |
| **prev_txs**     | []string | 输入的交易原文列表 |
| **txids**        | []string | 输入的交易id列表   |
| **input_indexs** | []uint32 | 输入的交易索引列表 |
| **addresses**    | []string | 输出的地址列表     |
| **amounts**      | []uint64 | 输出的数量列表     |
| **Return**       | string   | 初始的交易原文     |

#### **返回错误**

- `txids and indexs must be equal in length`
- `addresses and amounts must be equal in length`
- `Input count must be greater than 0`
- `Output count must be greater than 0`
- `Invalid Transaction`
- `Invalid Tx Input`
- `Invalid Tx Output`

---

### GetSighash(tx, txid, input_index, agg_pubkey, sigversion)**

#### **说明**

计算交易哈希(sighash)。一笔交易有多个输入，每个输入都需计算一个sighash,然后对该sighash进行签名得到signature。

#### **参数和返回值**

| **Name**        | **Type** | **Description**                                              |
| --------------- | -------- | ------------------------------------------------------------ |
| **tx**          | string   | generateRawTx返回的结果                                      |
| **txid**        | string   | 输入的交易id                                                 |
| **input_index** | uint32   | 输入的交易索引                                               |
| **agg_pubkey**  | string   | 输入是非门限地址时，填入""；门限地址时填入聚合公钥(getAggPublicKey) |
| **sigversion**  | uint32   | 输入是非门限地址时，填入0；输入是门限地址时，填入1；         |
| **Return**      | string   | 当前输入的交易哈希                                           |

#### **返回错误**

- `Compute Sighash Fail`

---

### GenerateSchnorrSignature(message, privkey)**

#### **说明**

非门限地址时，利用上述sighash和该函数计算签名

#### **参数和返回值**

| **Name**    | **Type** | **Description**                       |
| ----------- | -------- | ------------------------------------- |
| **message** | String   | 待签名的消息，即上面计算出来的sighash |
| **privkey** | String   | 签名者的私钥                          |
| **Return**  | String   | Schnorr签名                           |

#### **返回错误**

- `Invalid Signature`

---

### GetUnsignedTx(tx)**

#### **说明**

从`generateRawTx`生成的未签名的交易原文，携带有自定义的附加信息，不是有效的交易原文。`getUnsignedTx`的目的是生成有效的未签名的交易原文，能被BTC网络解析。

#### **参数和返回值**

| **Name**   | **Type** | **Description**                |
| ---------- | -------- | ------------------------------ |
| **tx**     | String   | 携带附加信息的未签名的交易原文 |
| **Return** | String   | 生成有效的未签名的交易原文     |

#### **返回错误**

- `Invalid Transaction`

---


### BuildTaprootTx(tx, signature, input_index)**

#### **说明**

非门限地址时，利用该函数将`generateSchnorrSignature`生成的签名组装进`generateRawTx`生成的原始交易。每一个输入都要进行一次签名，因此多个输入要组装多次。

#### **参数和返回值**

| **Name**        | **Type** | **Description**               |
| --------------- | -------- | ----------------------------- |
| **tx**          | string   | generateRawTx计算出的原始交易 |
| **signature**   | string   | 单个Schnorr签名               |
| **input_index** | uint32   | 输入的交易索引                |
| **Return**      | string   | 返回组装后的交易              |

#### **返回错误**

- *`Construct Tx Fail`*

---

### BuildThresholdTx(tx, agg_signature, agg_pubkey, control, input_index)**

#### **说明**

门限地址时，利用该函数将`Musig2`生成的聚合签名组装进`generateRawTx`生成的原始交易。每一个输入都要进行一次签名，因此多个输入要组装多次。

#### **参数和返回值**

| **Name**          | **Type** | **Description**               |
| ----------------- | -------- | ----------------------------- |
| **tx**            | string   | generateRawTx计算出的原始交易 |
| **agg_signature** | string   | Musig2聚合签名                |
| **agg_pubkey**    | string   | Musig2聚合公钥                |
| **control**       | string   | Mast生成的proof               |
| **input_index**   | uint32   | 输入的交易索引                |
| **Return**        | string   | 返回组装后的交易              |

#### **返回错误**

- *`Construct Tx Fail`*

---

### GetScriptPubkey(addr)**

#### **说明**

利用地址生成scirpt_pubkey,支持所有的地址格式。

#### **参数和返回值**

| **Name**   | **Type** | **Description** |
| ---------- | -------- | --------------- |
| **addr**   | String   | 地址            |
| **Return** | String   | scirpt_pubkey   |

#### **返回错误**

- `Invalid Address`

---

### GenerateSpentOutputs(prev_txs, input_indexs)**

#### **说明**

生成spend outputs。使用在Chainx的`createTaprootWithdrawTx`.

#### **参数和返回值**

| **Name**         | **Type** | **Description**       |
| ---------------- | -------- | --------------------- |
| **prev_txs**     | []string | 输入交易数组          |
| **input_indexs** | []uint32 | 输入交易索引数组      |
| **Return**       | string   | 序列化的spend outputs |

#### **返回错误**

- *`Invalid Spent Outputs`*

---

### GetMyAddress(pubkey, network)

#### **说明**

生成地址

#### **参数和返回值**

| **Name**    | **Type** | **Description**                                              |
| ----------- | -------- | ------------------------------------------------------------ |
| **pubkey**  | String   | 公钥                                                         |
| **network** | String   | 比特币网络类型，支持“mainnet”，“signet”， “testnet”， “regtest” |
| **Return**  | String   | 地址                                                         |

#### **返回错误**

- `Invalid Public Bytes`

---

## **Musig2**

下面是聚合签名和聚合公钥相关的函数

### GetMyPrivkey(phrase, pd_passphrase)**

#### **说明**

通过助记词和密码生成私钥

#### **参数和返回值**

| **Name**          | **Type** | **Description** |
| ----------------- | -------- | --------------- |
| **phrase**        | String   | 助记词          |
| **pd_passphrase** | String   | 密码            |
| **Return**        | String   | 私钥            |

#### **返回错误**

- `Construct Secret Key`

---

### GetMyPubkey(private)**

#### **说明**

通过私钥生成公钥

#### **参数和返回值**

| **Name**    | **Type** | **Description** |
| ----------- | -------- | --------------- |
| **private** | String   | 私钥            |
| **Return**  | String   | 公钥            |

#### **返回错误**

- *`Null KeyPair Pointer`*
- *`Normal Error`*

---

### GetRound1State()**

#### **说明**

Musig2生成第一轮的状态.

#### **参数和返回值**

| **Name**   | **Type** | **Description** |
| ---------- | -------- | --------------- |
| **Return** | *C.State | 第一轮状态      |

#### **返回错误**

- `null pointer`

---

### **GetRound1Msg(state)**

#### **说明**

通过第一轮状态生成消息，用于传递给其他参与者

#### **参数和返回值**

| **Name**   | **Type** | **Description** |
| ---------- | -------- | --------------- |
| **state**  | *C.State | 第一轮状态      |
| **Return** | string   | 第一轮消息      |

#### **返回错误**

- `Null Round1 State Pointer`
- *`Normal Error`*

---

### EncodeRound1State(state)**

#### **说明**

对第一轮状态序列化

#### **参数和返回值**

| **Name**   | **Type** | **Description** |
| ---------- | -------- | --------------- |
| **state**  | *C.State | 第一轮状态      |
| **Return** | string   | 序列化结果      |

#### **返回错误**

- `Null Round1 State Pointer`
- `Encode Fail`

---

### **DecodeRound1State(round1_state)**

#### **说明**

对第一轮状态反序列化

#### **参数和返回值**

| **Name**         | **Type** | **Description**           |
| ---------------- | -------- | ------------------------- |
| **round1_state** | string   | EncodeRound1State的输出值 |
| **Return**       | *C.State | 第一轮状态                |

#### **返回错误**

- `null pointer`

---

### **GetRound2Msg(state, msg, priv, pubkeys, received_round1_msg)**

#### **说明**

生成第二轮消息

#### **参数和返回值**

| **Name**                | **Type** | **Description**                        |
| ----------------------- | -------- | -------------------------------------- |
| **state**               | *C.State | DecodeRound1State的输出值              |
| **msg**                 | string   | 待签名的消息，通常是getSighash的返回值 |
| **priv**                | string   | 当前参与者私钥                         |
| **pubkeys**             | []string | 所有多签参与者公钥                     |
| **received_round1_msg** | []string | 接收到的其他多签参与者的第一轮消息     |
| **Return**              | string   | 第二轮消息                             |

#### **返回错误**

- `null pointer`

---

### GetAggSignature(round2_msg)**

#### **说明**

返回聚合签名的结果

#### **参数和返回值**

| **Name**       | **Type** | **Description**        |
| -------------- | -------- | ---------------------- |
| **round2_msg** | String   | 所有参与者的第二轮消息 |
| **Return**     | String   | 签名结果               |

#### **返回错误**

- `Normal Error`
- `Null Round2 State Pointer`

---

### **GetAggPublicKey(pubkeys)**

#### **说明**

生成聚合公钥

#### **参数和返回值**

| **Name**    | **Type** | **Description**  |
| ----------- | -------- | ---------------- |
| **pubkeys** | String[] | 待聚合的公钥列表 |
| **Return**  | String   | 聚合公钥         |

#### **返回错误**

- `Normal Error`

---

## **Mast**

下面是生成门限地址和proof相关的函数

### GenerateThresholdPubkey(pubkeys, threshold)

#### **说明**

生成门限公钥

#### **参数和返回值**

| **Name**      | **Type** | **Description** |
| ------------- | -------- | --------------- |
| **pubkeys**   | []string | 所有的公钥列表  |
| **threshold** | uint8    | 阈值            |
| **Return**    | string   | 聚合公钥        |

#### **返回错误**

- `Invalid Public Bytes`

---

### GenerateControlBlock(pubkeys, threshold, aggPubkey)**

#### **说明**

生成proof

#### **参数和返回值**

| **Name**      | **Type** | **Description**          |
| ------------- | -------- | ------------------------ |
| **pubkeys**   | []string | 所有的公钥列表           |
| **threshold** | uint8    | 阈值                     |
| **aggPubkey** | string   | 本次多签参与者的聚合公钥 |
| **Return**    | string   | proof                    |

#### **返回错误**

- `Invalid Public Bytes`

---

# **Example**

下面示例提供了：构造非门限地址，非门限地址的花费，构造门限签名地址，门限签名地址花费。完整代码可以在[demo](./demo/musig2demo.go)中查看。

## **Details**

### 生成非门限签名地址

1. 传入助记词和密码，生成私钥

   ~~~go
   private0, err := musig2.GetMyPrivkey(PHRASE0, "")
   ~~~

2. 生成公钥

   ~~~go
   pubkey0, err := musig2.GetMyPubkey(private0)
   ~~~

3. 生成地址

   ~~~go
   address0, err := musig2.GetMyAddress(pubkey0, "signet")
   ~~~

### 非门限签名地址的花费

1. **通过`generateRawTx`创建一笔未签名的交易**。txids和indexs用于构造交易的所有输入，一个txid和一个index用来定位唯一一笔未花费的输出。下面**prev_txs,txids和input_indexs长度一致并且一一对应**。addresses和amounts用于构造交易的所有输出，一个adddress和一个amount表示向一个地址发送多少币。**adddress没有顺序要求，只需amounts一一对应即可**。与这里`1f8e0f7dfa37b184244d022cdf2bc7b8e0bac8b52143ea786fa3f7bbe049eeae`和`1`唯一确定了一笔未花费的输出，这个未花费的输出所属的地址是一个**非门限地址**。用txid可以查询到相应的p rev_tx。`35516a706f3772516e7751657479736167477a6334526a376f737758534c6d4d7141754332416255364c464646476a38`代表着`op_return`,它所对应的amout为0。`tb1pn202yeugfa25nssxk2hv902kmxrnp7g9xt487u256n20jgahuwasdcjfdw`是接收方的地址，`100000`是转账金额。`tb1pexff2s7l58sthpyfrtx500ax234stcnt0gz2lr4kwe0ue95a2e0srxsc68`是找零地址，`400000`是找零金额。计算方式参考**手续费和找零余额计算**。

   ~~~go
   prev_txs :=[]string{"020000000001014be640313b023c3c731b7e89c3f97bebcebf9772ea2f7747e5604f4483a447b601000000000000000002a0860100000000002251209a9ea267884f5549c206b2aec2bd56d98730f90532ea7f7154d4d4f923b7e3bbc027090000000000225120c9929543dfa1e0bb84891acd47bfa6546b05e26b7a04af8eb6765fcc969d565f01404dc68b31efc1468f84db7e9716a84c19bbc53c2d252fd1d72fa6469e860a74486b0990332b69718dbcb5acad9d48634d23ee9c215ab15fb16f4732bed1770fdf00000000"}
   txids := []string{"1f8e0f7dfa37b184244d022cdf2bc7b8e0bac8b52143ea786fa3f7bbe049eeae"}
   input_indexs := []uint32{1}
   addresses := []string{"tb1pn202yeugfa25nssxk2hv902kmxrnp7g9xt487u256n20jgahuwasdcjfdw", "35516a706f3772516e7751657479736167477a6334526a376f737758534c6d4d7141754332416255364c464646476a38", "tb1pexff2s7l58sthpyfrtx500ax234stcnt0gz2lr4kwe0ue95a2e0srxsc68"}
   amounts = []uint64{100000, 0, 400000}
   base_tx, err := musig2.GenerateRawTx(prev_txs, txids, input_indexs, addresses, amounts)
   ~~~
   
2. **对要花费的输出进行签名**。对要花费的UTXO进行签名首先要计算出这笔未花费输出的sighash，签名是对sighash进行签名.

   txid以及input_index用来定位那笔要花费的输出，agg_pubkey对于非门限签名地址填空字符串`""`，sigversion对于非门限签名地址填0，tx是当前构造的交易
   。**注意计算sighash的时候，永远要用上面`generateRawTx`构造出的结果不能改变。**
   ~~~go
   sighash, err := musig2.GetSighash(baseTx, txids[i], inputIndexs[i], "", 0)
   ~~~

   计算完sighash后，再使用私钥对其进行签名。message就是指sighash，privkey就是私钥。

   ~~~go
   schnorrSignature, err := musig2.GenerateSchnorrSignature(sighash, privkey)
   ~~~

3. **将上面的签名组装进交易**。tx就是当前要构造的交易，txid和input_index仍然用来定位tx中签名对应的输入。

   ~~~go
   taprootTx, err = musig2.BuildTaprootTx(baseTx, schnorrSignature, txids[i], inputIndexs[i])
   ~~~
   
   **注意如果tx中有多个输入，那么需要重复Step2和Step3对每个输出进行签名并添加到tx中，如下图所示的for循环：。**
   
   ![](https://cdn.jsdelivr.net/gh/hacpy/PictureBed@master/Document/16523346050391652334604157.png)

### 生成门限签名地址

1. 如下生成一个2-of-3的门限签名地址,。首先传入所有参与者的公钥和阈值即可生成门限公钥。

   ~~~go
   thresholdPubkey, err := musig2.GenerateThresholdPubkey([]string{pubkey0, pubkey1, pubkey2}, 2)
   ~~~

2. 再将公钥编码成地址，就可以得到门限地址

   ~~~go
   thresholdAddress, err := musig2.GetMyAddress(thresholdPubkey, "signet")
   ~~~

### 门限签名地址的花费

1. **通过`generateRawTx`创建一笔未签名的交易**。txids和indexs用于构造交易的所有输入，一个txid和一个index用来定位唯一一笔未花费的输出。下面**prev_txs,txids和input_indexs长度一致并且一一对应**。addresses和amounts用于构造交易的所有输出，一个adddress和一个amount表示向一个地址发送多少币。**adddress没有顺序要求，只需amounts一一对应即可**。这里`8e5d37c768acc4f3e794a10ad27bf0256237c80c22fa67117e3e3e1aec22ea5f`和`0`唯一确定了一笔未花费的输出，注意这个未花费的输出所属的地址是一个**门限地址**。用txid可以查询到相应的p rev_tx。`tb1pexff2s7l58sthpyfrtx500ax234stcnt0gz2lr4kwe0ue95a2e0srxsc68`是接收方的地址，`50000`是转账金额。`tb1pn202yeugfa25nssxk2hv902kmxrnp7g9xt487u256n20jgahuwasdcjfdw`是找零地址，`40000`是找零金额。当然这里也可以带`op_return`。计算方式参考**手续费和找零余额计算**

   ~~~go
   prev_txs = []string{"02000000000101aeee49e0bbf7a36f78ea4321b5c8bae0b8c72bdf2c024d2484b137fa7d0f8e1f01000000000000000003a0860100000000002251209a9ea267884f5549c206b2aec2bd56d98730f90532ea7f7154d4d4f923b7e3bb0000000000000000326a3035516a706f3772516e7751657479736167477a6334526a376f737758534c6d4d7141754332416255364c464646476a38801a060000000000225120c9929543dfa1e0bb84891acd47bfa6546b05e26b7a04af8eb6765fcc969d565f01409e325889515ed47099fdd7098e6fafdc880b21456d3f368457de923f4229286e34cef68816348a0581ae5885ede248a35ac4b09da61a7b9b90f34c200872d2e300000000"};
   txids = []string{"8e5d37c768acc4f3e794a10ad27bf0256237c80c22fa67117e3e3e1aec22ea5f"};
   input_indexs = []uint32{0};
   addresses = []string{"tb1pexff2s7l58sthpyfrtx500ax234stcnt0gz2lr4kwe0ue95a2e0srxsc68","tb1pn202yeugfa25nssxk2hv902kmxrnp7g9xt487u256n20jgahuwasdcjfdw"};
   amounts = []uint64{50000, 40000};
   string base_tx = musig2.GenerateRawTx(prev_txs, txids, input_indexs, addresses, amounts);
   ~~~
   
2. **对要花费的输出进行签名**。对要花费的UTXO进行签名首先要计算出这笔未花费输出的sighash，签名是对sighash进行签名。

   txid以及input_index用来定位那笔要花费的输出，agg_pubkey对于门限签名地址填空字符串聚合公钥，如下是B和C两个人进行聚签花费，那么就填入B和C的聚合公钥。sigversion对于门限签名地址填1，tx是当前构造的交易。

   **计算sighash**

   ~~~go
   pubkeyBC, err := musig2.GetAggPublicKey([]string{pubkeyB, pubkeyC})
   sighash, err := musig2.GetSighash(baseTx, txids[i], inputIndexs[i], pubkeyBC, 1)
   ~~~

   **计算签名**：计算完sighash后,B和C两个人利用Musig2进行聚合签名。签名的消息就是sighash。

   ~~~go
   round1State0 := musig2.GetRound1State()
   round1State0Encode, err := musig2.EncodeRound1State(round1State0)
   round1State0 = musig2.DecodeRound1State(round1State0Encode)
   round1State1 := musig2.GetRound1State()
   
   round1Msg0, err := musig2.GetRound1Msg(round1State0)
   round1Msg1, err := musig2.GetRound1Msg(round1State1)
   round2Msg0, err := musig2.GetRound2Msg(round1State0, sighash, privateB, []string{pubkeyB, pubkeyC}, []string{round1Msg1})
   round2Msg1, err := musig2.GetRound2Msg(round1State1, sighash, privateC, []string{pubkeyB, pubkeyC}, []string{round1Msg0})
   multiSignature, err := musig2.GetAggSignature([]string{round2Msg0, round2Msg1})
   controlBlock, err := musig2.GenerateControlBlock([]string{pubkeyA, pubkeyB, pubkeyC}, 2, pubkeyBC)
   ~~~

   ​	下面是对上述Musig2多签的过程的详细介绍，分为如下几步：

   1. 生成第一轮的状态

      ~~~go
      round1State0 := musig2.GetRound1State()
      ~~~

   2. 通过第一轮状态获取第一轮消息,并传递给其他签名参与者。

      ~~~go
      round1Msg0, err := musig2.GetRound1Msg(round1State0)
      ~~~

   3. 拿到其他签名参与者的第一轮消息，生成第二轮消息，并传递给其他参与者。`received_round1_msg`是接收到的其他参与者的第一轮消息。`pubkeys`是所有参与者的公钥。`msg`是待签名的消息。`state`是第一轮的状态。`priv`是签名者私钥。

      ~~~go
      round2Msg0, err := musig2.GetRound2Msg(round1State0, sighash, privateB, []string{pubkeyB, pubkeyC}, []string{round1Msg1})
      ~~~

   4. 利用所有参与者的第二轮消息，生成聚合签名。`round2_msg` 是所有参与者的第二轮消息。

      ~~~go
      string multi_signature = getAggSignature(new []string{round2_msg0, round2_msg1})
      ~~~

   **计算proof**: 门限签名的花费不仅需要签名，还要计算proof。需要传入所有人的公钥，阈值和本次签名参与者B和C的聚合公钥。

   ~~~java
   controlBlock, err := musig2.GenerateControlBlock([]string{pubkeyA, pubkeyB, pubkeyC}, 2, pubkeyBC)
   ~~~

3. **将上面的签名和proof组装进行交易**。tx就是当前要构造的交易，agg_signature是B和C的聚合签名，agg_pubkey是B和C的聚合公钥，txid和input_index仍然用来定位tx中签名对应的输入，txid和input_index对应的未花费输出与第二步是对应的。

   ~~~go
   thresholdTx, err = musig2.BuildThresholdTx(baseTx, multiSignature, pubkeyBC, controlBlock, txids[i], inputIndexs[i])
   ~~~
   
   **注意如果tx中有多个输入，那么需要重复Step2和Step3对每个输出进行签名并添加到tx中，如下图所示的for循环：**

   ![](https://cdn.jsdelivr.net/gh/hacpy/PictureBed@master/Document/16523351090411652335108767.png)

## 手续费和找零余额计算

背景: A要转账给`B 2BTC`,` C 3BTC`

1. 通过A的地址找到所有未花费的交易txids和余额，并从大到小排序，假设为`[(txid1, 4), (txid2, 2), (tixd3, 1), (tixd4, 1)]`。

2. 对txids和余额列表累加并找到大于输出金额2+3=5的txid，也就是txid2，未找到则返回不允许转账。

3. 从txid2向后顺延一位，用`[(txid1, 4), (txid2, 2), (tixd3, 1)]`作为输入。如果txid2是最后一个，用`[(txid1, 4), (txid2, 2)]`作为输入。

4. 利用输入和输出的个数以及如下公式，估计交易字节数：

   **非门限地址花费的字节数估计**

   ~~~
   105 + 58 * input_count(threshold_address) + 43 * output_count
   ~~~
   `input_count(taproot_address)`表示非门限地址花费时输入txid的个数

   **门限地址的字节数估计**

   ```
   105 + 141 * input_count(threshold_address) + 43 * output_count
   ```
   `input_count(threshold_address)`表示门限地址花费时输入txid的个数

5. 利用字节数乘以当前`FEE RATES`得到交易手续费。

6. 将`输入总金额 - （输出总金额+手续费）`得到`找零金额`。 如果为负则没有找零（即输出列表不填入找零地址和金额），此时交易手续费成了`输入总金额 - 输出总金额`。

# Other

- 将地址转成比特币交易中输出的锁定脚本script_pubkey

  ~~~go
  scriptPubkey, err := musig2.GetScriptPubkey("tb1pn202yeugfa25nssxk2hv902kmxrnp7g9xt487u256n20jgahuwasdcjfdw")
  ~~~



- 从`generateRawTx`生成的未签名的交易原文中提取有效的未签名的交易原文。

  ~~~go
  unsignedTx, err := musig2.GetUnsignedTx(baseTx)
  ~~~
