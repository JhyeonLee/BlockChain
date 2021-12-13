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