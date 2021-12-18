Make BlockChain

Why?

For building deep learning model with blockchain together.

What if?

Each layer with blockchain : perceptron == block and interconnection of block perceptrons

or

Each model with blockchain : model == block and interconnection of block models

On this Repository
1. Basic Blockchain Concept
2. Explorer with *.gohtml
3. REST API ~> Router of [Gorilla MUX](https://github.com/gorilla/mux)
4. Simple CLI, and It will be by [cobra CLI](https://github.com/spf13/cobra)
5. Database of [bolt](https://github.com/boltdb/bolt) ~> hash(key): block(value)

    - for checking db : [boltbrowser](https://github.com/br0xen/boltbrowser) and [boltdbweb](https://github.com/evnix/boltdbweb)
6. PoW(Proof of Work) about Mining ~> ex. verified block's first 19 digits of hash are 0

    - hard to solve but easy to verify
7. Transaction & uTxO(Unspent Transaction Outputs)

    - Tx : transaction
    - TxIn : money before transaction
    - TxOut : money after transaction

    - Mempool(Memory Pool) == Unconfirmed
    - Check whether is Unsent Transaction
    - Check whether trasaction is on Mempool

    - Refactoring
        - Method: should mutate struct ~>ex. `func (b *blockchain) AddBlock()`
        - if not, it is function ~> ex. `func Blocks(b *blockchain) []*Block`
    
    - Deadlock
        - Because no call to Do returns until the one call to f returns, if f causes Do to be called, it will deadlock.
8. Wallet

    - check the owner owns unspent transaction output ~> hash and signature
    - verify the owner approves the transaction

    - how signature and verification work ~> public key and private key by elliptic curve cryptography
    - persistency for wallet ~> backend for wallet
    - implement of signature and verification, applied to transaction

    - how to verify someone is unspent transaction ouput's owner ~> signature from private key and address from public key / TxIn to know unspent TxOut
    - TxIn has signature, and signature is by private key
    - TxOut has address, address is where you sent, and address is public key
    - Tx{ TxIn[ (TxOut1) (TxOut2) ], Sign } ~~~>>> TxIn.Sign + TxOut1.Address = true / false
9. P2P

    - P2P Network : Network, many nodes running same software ~> Decentralization
    - WebSocket : Notification system in real time, nodes discovering and conneting each other ~> making mesh and updating each other, broadcasting transaction to whole network
        - protocol like http, but difference is http is stateless(after send and receive request, no memory of connection between user and server ~> after request, server do not remember user)
        - websocket is alive connection(bi-directional connection) ~> alive bi-directional connection with all nodes ~> upgrade http(stateless protocol, server forger user) to websocket(statefull protocol, server remember user)
    - Channel : Connecting peers for network ~> occur race codition problem, beacuse of concurrnecy
    - Race Condition Problem ~>  new version of bolt [bbolt](https://github.com/etcd-io/bbolt) & Mutexes fix it
        - data race happens when two goRoutine access same data in an application
        - one goRoutine reads the data and the other updates the data
        - solution: locking the data ~> mutex
    - Broadcasting all information with all nodes