# 两阶段事务业务合约

## 介绍
fundcc是两阶段事务业务合约，合约包含资源的创建查询和基金购买等方法。本合约基于Fabric 2.0语法开发，可同时兼容fabric 2.x网络和fabric 1.4网络。本合已合并yttmcc（旧跨链合约）并适配crosscc-fabric（新跨链合约）合约。

## 合约更新步骤
#### 1.打包新的链码

```shell
peer lifecycle chaincode package demo_1.0.tar.gz \
--path /home/centos/go/src/github.com/chaincode/go/demo \
--lang golang \
--label demo_1.0
```

#### 2.切换到Org1

```shell
export CORE_PEER_LOCALMSPID="Org1MSP"
export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt
export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp
export CORE_PEER_ADDRESS=localhost:7051
```

#### 3.安装链码Org1

```shell
peer lifecycle chaincode install demo_1.0.tar.gz
```

#### 4.查询安装情况

```shell
peer lifecycle chaincode queryinstalled
#Installed chaincodes on peer:
#Package ID: basic_1.0:dee2d612e15f5059478b9048fa4b3c9f792096554841d642b9b59099fa0e04a4, Label: basic_1.0
#Package ID: basic_2.0:aa2b26ca109e0fe1ba12782dd8683769e2a2b70d05027005199202548ef2b800, Label: basic_2.0
```

#### 5.设置package_id

```shell
export NEW_CC_PACKAGE_ID=basic_3.0:afd55fbf1bd34ead75325c541dec705d92418348da02a623f3a445f7ee02aad9
#demo
export DM_PACKAGE_ID=demo_1.0:1c3bdf6b2d7388beb1dacecb360a62d24d2db163198b008eb78575c4bef9af43
```

#### 6.approve

```shell
peer lifecycle chaincode approveformyorg -o localhost:7050 \
--ordererTLSHostnameOverride orderer.example.com \
--channelID mychannel \
--name demo \
--version 1.0 \
--package-id $DM_PACKAGE_ID \
--sequence 1 \
--init-required \
--tls \
--cafile "${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem"
```

#### 7.切换到Org2

```shell
export CORE_PEER_LOCALMSPID="Org2MSP"
export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt
export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/org2.example.com/users/Admin@org2.example.com/msp
export CORE_PEER_ADDRESS=localhost:9051
```

#### 8.安装链码Org2

```shell
// 参数同上
```

#### 9.approve

```shell
// 参数同上
```

#### 10.check commit readiness

```shell
peer lifecycle chaincode checkcommitreadiness \
--channelID mychannel \
--name demo \
--version 1.0 \
--sequence 1 \
--tls \
--init-required \
--cafile "${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem" --output json
```

#### 11.commit

```shell
peer lifecycle chaincode commit -o localhost:7050 \
--ordererTLSHostnameOverride orderer.example.com \
--channelID mychannel \
--name demo \
--version 1.0 \
--sequence 1 \
--init-required \
--tls \
--cafile "${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem" \
--peerAddresses localhost:7051 \
--tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt" \
--peerAddresses localhost:9051 \
--tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt"
```

#### 12.docker ps

```shell
docker ps
CONTAINER ID   IMAGE                                                                                                                                                                    COMMAND                  CREATED          STATUS          PORTS                                                                                                                             NAMES
c4fcdc1a2964   dev-peer0.org1.example.com-basic_3.0-afd55fbf1bd34ead75325c541dec705d92418348da02a623f3a445f7ee02aad9-ed525edb9b6fa67c1c41e7493c4b0498a8b3085d89ad96ec10191a21c72eacda   "chaincode -peer.add…"   4 seconds ago    Up 3 seconds                                                                                                                                      dev-peer0.org1.example.com-basic_3.0-afd55fbf1bd34ead75325c541dec705d92418348da02a623f3a445f7ee02aad9
77babe521eb0   dev-peer0.org2.example.com-basic_3.0-afd55fbf1bd34ead75325c541dec705d92418348da02a623f3a445f7ee02aad9-bc157a4709e36b6f85e85f85ea83a28de1caaa3fd1e93b035562bbb9118e69c1   "chaincode -peer.add…"   5 seconds ago    Up 3 seconds
```

#### 13.invoke

```shell
peer chaincode invoke -o localhost:7050 \
--ordererTLSHostnameOverride orderer.example.com \
--tls \
--cafile "${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem" \
-C mychannel \
-n demo \
--peerAddresses localhost:7051 \
--tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt" \
--peerAddresses localhost:9051 \
--tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt" \
-c '{"function":"init","Args":[]}'
```

#### 14.query

```shell
peer chaincode invoke -o localhost:7050 \
--ordererTLSHostnameOverride orderer.example.com \
--tls \
--cafile "${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem" \
-C mychannel \
-n demo \
--peerAddresses localhost:7051 \
--tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt" \
--peerAddresses localhost:9051 \
--tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt" \
-c '{"function":"list_asset","Args":[]}'
```

