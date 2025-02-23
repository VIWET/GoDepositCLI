# Go Deposit CLI

CLI for generating of Bahamut chain and Ethereum 2.0 keys written in Golang

## Build

### Requirements

- make

- go v1.23.3

### Build for Bahamut chain

```Bash
git clone https://github.com/VIWET/GoDepositCLI.git && cd GoDepositCLI
make
```

### Build for Ethereum

```Bash
git clone https://github.com/VIWET/GoDepositCLI.git && cd GoDepositCLI
make NETWORK=ethereum
```

## TODO

- [ ] Config form terminal UI

- [ ] Colorschemes

- [ ] Mnemonic input refactoring

- [ ] Password generating

- [ ] Execution and Consensus layer clients

## Deposit generating

![Deposits generating](images/demo.gif)

#### Flags

- `--start-index`
	
    Aliases: `--index`,  `--start`, `--from`, `-i`

	Default: `0`

    The starting key index for generating deposits, keystores, or BLS to Execution messages. For example, if you previously generated `N` keys using a mnemonic, use `--start-index=N` to generate the next batch.

- `--number`

	Aliases: `--num value`, `-n`

	Default: `1`

    The number of deposits and keystores, or BLS to Execution messages to generate, starting from `--start-index`.

- `--amounts`

	Aliases: `--amount`, `-a`

    The amount of Ether to deposit. You can specify amounts per key index using `--amount=<key_index>:<amount>` or set a default amount for all deposits with `--amount=<amount>`.
    
    **Examples:**  

    - Single-line format (comma-separated values):

        ```Bash
        --amounts="<amount_1>,<key_index_2>:<amount_2>,..."
        ```
    - Multi-line format:

        ```Bash
        --amount="<amount_1>" \
        --amount="<key_index_2>:<amount_2>"
        ```

    Amount values support `GWEI`, `FTN`, or `ETH` prefixes. If no prefix is provided, the default is `FTN`/`ETH`.

- `--withdrawal-addresses`

	Aliases: `--withdrawal-address`, `--withdraw-to`, `--w`

    Withdrawal addresses for deposits. Similar to `--amounts`, you can specify a unique address per key index (`--withdraw-to="<key_index>:<address>"`) or set a default withdrawal address (`--withdraw-to=<address>`). Addresses must be 20-byte hex strings.

- `--contracts` (only for Bahamut chain)

	Aliases: `--contract`, `-c`

    Contract addresses used for validator activity calculations. Unlike `--amounts` and `--withdrawal-addresses`, this flag does not support a default value unless `--number=1`. The format is `--contract="<key_index>:<contract>"`. Contract values must be 20-byte hex strings.

- `--validator-indices`

    Aliases: `--validator-index`, `--indices`, `--vi`

    Specifies validator indices in the format `--validator-index="<key_index>:<validator_index>"`. All indices must be explicitly defined, except when `--number=1`.

- `--directory`

	Aliases: `--dir`, `-d`

	Default: `./validators_data`

	Directory where all generated data will be stored.

- `--keystore-kdf`

    The key derivation function used for keystore generation. Supported values: `scrypt`, `pbkdf2`.

- `--chain`

	Aliases: `--network`

	Default: `mainnet`

    The blockchain network configuration used for signing domain calculations.  

    - **Bahamut:** Supports `mainnet`, `sahara`, and `horizon`.  

    - **Ethereum:** Supports `mainnet` and `holesky`.

    If specifying an unknown network, `--genesis-fork` and `--genesis-validators-root` must be provided.

- `--genesis-fork`

    A unique 4-byte hex string used for signing domain calculations.

- `--genesis-validators-root`

    A unique 32-byte hex string used for signing domain calculations.

- `--password`

    The password used for keystore encryption. If not set, `staking-cli` will prompt you to enter it.  

    **Note:** If using `--non-interactive`, this flag is required.

- `--config`

	Aliases: `--cfg`

    Path to a configuration file for generating deposits and keystores.

	Example:

    - **Deposit config**

        ```json
        {
            "start_index": 0,
            "number": 4,
            "chain_config": {
                "name": "devnet",
                "genesis_fork_version": "0x00004242",
                "genesis_validators_root": "0x4242424242424242424242424242424242424242424242424242424242424242"
            },
            "mnemonic_config": {
                "language": "english",
                "bitlen": 256
            },
            "amounts": {
                "default": 8192000000000,
                "config": {
                    "0": 256000000000,
                    "1": 4096000000000
                }
            },
            "contract_addresses": {
                "config": {
                    "0": "0xdeaddeaddeaddeaddeaddeaddeaddeaddead0001",
                    "1": "0xdeaddeaddeaddeaddeaddeaddeaddeaddead0002"
                }
            },
            "withdrawal_addresses": {
                "default": "0xdeaddeaddeaddeaddeaddeaddeaddeaddeaddead",
                "config": {
                    "0": "0xdeaddeaddeaddeaddeaddeaddeaddeaddead0003",
                    "1": "0xdeaddeaddeaddeaddeaddeaddeaddeaddead0004"
                }
            },
            "directory": "./my_keys",
            "kdf": "scrypt"
        }
        ```

    - **BLS To Execution config**

        ```json
        {
            "start_index": 0,
            "number": 3,
            "chain_config": {
                "name": "devnet",
                "genesis_fork_version": "0x00004242",
                "genesis_validators_root": "0x4242424242424242424242424242424242424242424242424242424242424242"
            },
            "mnemonic_config": {
                "language": "english",
                "bitlen": 256
            },
            "validator_indices": {
                "config": {
                    "0": 0,
                    "1": 1,
                    "2": 2
                }
            },
            "withdrawal_addresses": {
                "default": "0xdeaddeaddeaddeaddeaddeaddeaddeaddeaddead",
                "config": {
                    "0": "0xdeaddeaddeaddeaddeaddeaddeaddeaddead0001",
                    "1": "0xdeaddeaddeaddeaddeaddeaddeaddeaddead0002"
                }
            },
            "directory": "./my_bls_to_execution"
        }
	    ```

- `--mnemonic-bitlen`

	Aliases: `--strength`,  `--bitlen`, `--bl`, `-s`

	Default: `256`

    The bit strength of the seed phrase. Options: `128`, `160`, `192`, `224`, `256`. This affects mnemonic length (`12`, `15`, `18`, `21`, or `24` words).

- `--mnemonic-language`

	Aliases: `--language`, `--lang`, `-l`

	Defulat: `english`

    The language for the mnemonic phrase. Supported languages: English, Chinese (Simplified & Traditional), Czech, French, Italian, Japanese, Korean, Portuguese, and Spanish.

- `--mnemonic`

    The seed phrase used for key generation. If not set, staking-cli will prompt you to enter it.

    **Note:** If using `--non-interactive`, this flag is required.

### Examples

#### New Mnemonic

```Bash
./bin/staking-cli deposit \
    --start-index=0 \
    --number=4 \
    --chain="devnet" \
    --genesis-fork="0x00004242" \
    --genesis-validators-root="0x4242424242424242424242424242424242424242424242424242424242424242" \
    --directory="./my_keys" \
    new-mnemonic \
    --amount="8192" \
    --amount="0:256FTN" \
    --amount="1:4096000000000GWEI" \
    --contract="0:0xdeaddeaddeaddeaddeaddeaddeaddeaddead0001" \
    --contract="1:0xdeaddeaddeaddeaddeaddeaddeaddeaddead0002" \
    --withdraw-to="0xdeaddeaddeaddeaddeaddeaddeaddeaddeaddead" \
    --withdraw-to="0:0xdeaddeaddeaddeaddeaddeaddeaddeaddead0003" \
    --withdraw-to="1:0xdeaddeaddeaddeaddeaddeaddeaddeaddead0004" \
    --bitlen=256 \
    --password="password" \
    --language="english"
```

#### Existing Mnemonic

```Bash
./bin/staking-cli deposit \
    --start-index=4 \
    --number=4 \
    --chain="devnet" \
    --genesis-fork="0x00004242" \
    --genesis-validators-root="0x4242424242424242424242424242424242424242424242424242424242424242" \
    --directory="./my_keys" \
    existing-mnemonic \
    --amount="8192" \
    --amount="4:256FTN" \
    --amount="5:4096000000000GWEI" \
    --contract="4:0xdeaddeaddeaddeaddeaddeaddeaddeaddead0005" \
    --contract="5:0xdeaddeaddeaddeaddeaddeaddeaddeaddead0006" \
    --withdraw-to="0xdeaddeaddeaddeaddeaddeaddeaddeaddeaddead" \
    --withdraw-to="4:0xdeaddeaddeaddeaddeaddeaddeaddeaddead0007" \
    --withdraw-to="5:0xdeaddeaddeaddeaddeaddeaddeaddeaddead0008" \
    --password="password" \
    --mnemonic="abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon about" \
    --language="english"
```


#### BLS To Execution

```bash
./bin/staking-cli bls-to-execution \
    --start-index=0 \
    --number=3 \
    --chain="devnet" \
    --genesis-fork="0x00004242" \
    --genesis-validators-root="0x4242424242424242424242424242424242424242424242424242424242424242" \
    --validator-index="0:0" \
    --validator-index="1:1" \
    --validator-index="2:2" \
    --withdraw-to="0xdeaddeaddeaddeaddeaddeaddeaddeaddeaddead" \
    --withdraw-to="0:0xdeaddeaddeaddeaddeaddeaddeaddeaddead0001" \
    --withdraw-to="1:0xdeaddeaddeaddeaddeaddeaddeaddeaddead0002" \
    --directory="./my_bls_to_execution" \
    --mnemonic="abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon about" \
    --language="english"
```
