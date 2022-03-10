# blockchain

TODO:
1. Need to iterate through keys of db to get the first block
    - Right now if a block chain has already been created, and then you create a block, trying
    to create a blockchain again will give you the most recent block created. We always want to
    return the genesis block
2. Iterate through keys to print out the block chain
