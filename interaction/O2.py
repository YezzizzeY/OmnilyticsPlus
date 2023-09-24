from web3 import Web3, AsyncWeb3

w3 = Web3(Web3.HTTPProvider('http://127.0.0.1:7545'))

address = '0xa8f8f562ec1DC16469C41ad9f8b65b318d45F5F7'
abi = '''[
	{
		"inputs": [
			{
				"internalType": "uint256",
				"name": "key",
				"type": "uint256"
			}
		],
		"name": "calculateAndStoreSquareSum",
		"outputs": [],
		"stateMutability": "nonpayable",
		"type": "function"
	},
	{
		"inputs": [],
		"name": "getCurrentKey",
		"outputs": [
			{
				"internalType": "uint256",
				"name": "",
				"type": "uint256"
			}
		],
		"stateMutability": "view",
		"type": "function"
	},
	{
		"inputs": [
			{
				"internalType": "uint256",
				"name": "key",
				"type": "uint256"
			}
		],
		"name": "getSquareSum",
		"outputs": [
			{
				"internalType": "int256",
				"name": "",
				"type": "int256"
			}
		],
		"stateMutability": "view",
		"type": "function"
	},
	{
		"inputs": [
			{
				"internalType": "int256[]",
				"name": "grad",
				"type": "int256[]"
			},
			{
				"internalType": "int256[]",
				"name": "beta",
				"type": "int256[]"
			},
			{
				"internalType": "int256[]",
				"name": "sigma",
				"type": "int256[]"
			}
		],
		"name": "storeData",
		"outputs": [],
		"stateMutability": "nonpayable",
		"type": "function"
	}
]'''
contract_instance = w3.eth.contract(address=address, abi=abi)



# ----------------------------------------------------------------------------------------------------------------------#

# ----------------------------------------------------------------------------------------------------------------------#


def read_and_convert_to_int256(filename):
    # Read the file
    with open(filename, "r") as file:
        content = file.read()

    # Split the content by spaces and newlines
    parts = content.split()

    # Try to convert each part to a float
    float_numbers = []
    for part in parts:
        try:
            num = float(part)
            float_numbers.append(num)
        except ValueError:
            pass

    # Convert to Solidity int256 format
    int256_numbers = [int(num * 10**18) for num in float_numbers]

    return int256_numbers

# File name parts
file_parts = ["fc1", "fc2", "fc3", "last"]

# Lists to store the numbers from each category
beta_list = []
grad_list = []
sigma_list = []

# Read and convert numbers from each file
for part in file_parts:
    for prefix in ["beta", "grad", "sigma"]:
        filename = f"mlp grad+beta+sigma/{prefix}{part}.txt"
        uint256_numbers = read_and_convert_to_int256(filename)
        
        # Append the numbers to the respective list
        if prefix == "beta":
            beta_list.extend(uint256_numbers)
        elif prefix == "grad":
            grad_list.extend(uint256_numbers)
        elif prefix == "sigma":
            sigma_list.extend(uint256_numbers)

# Print the total number of elements in each list
print(f"Total number of elements in beta_list: {len(beta_list)}")
print(f"Total number of elements in grad_list: {len(grad_list)}")
print(f"Total number of elements in sigma_list: {len(sigma_list)}")




# ----------------------------------------------------------------------------------------------------------------------#

# ----------------------------------------------------------------------------------------------------------------------#






def validate_int256_numbers(numbers):
    for num in numbers:
        assert isinstance(num, int), f"Number {num} is not an integer!"
        assert -2**255 <= num < 2**255, f"Number {num} is out of range for int256!"

def transact_data(key, start_idx, end_idx, grad, beta, sigma):
    validate_int256_numbers(grad)
    validate_int256_numbers(beta)
    validate_int256_numbers(sigma)

    print(f"Storing array indices: {start_idx} to {end_idx-1}")
    tx_hash = contract_instance.functions.storeData(grad, beta, sigma).transact({'from': '0xCeE4aecD7Af726D8338ef6e6eAdF12591E6e3b28', 'gas':999999999999999})

    tx_receipt = w3.eth.wait_for_transaction_receipt(tx_hash)
    assert tx_receipt['status'] == 1

    gas_used = tx_receipt['gasUsed']
    print(f"Data successfully stored! Gas used: {gas_used}")
    return gas_used

print(beta_list[0:50])
total_gas_used = 0

# Process each list in batches of 50
for i in range(0, len(grad_list), 2500):
    grad_batch = grad_list[i:i+2500]
    beta_batch = beta_list[i:i+2500]
    sigma_batch = sigma_list[i:i+2500]

    gas_used_for_batch = transact_data(i // 50 + 1, i, i+len(grad_batch), grad_batch, beta_batch, sigma_batch)
    total_gas_used += gas_used_for_batch

print(f"Total gas used for all transactions: {total_gas_used}")




# ----------------------------------------------------------------------------------------------------------------------#

# ----------------------------------------------------------------------------------------------------------------------#



# # 
# total_gas_used = 0  # To accumulate gas costs of all transactions

# # Iterate over key values from 1 to 149
# for key_to_calculate in range(1, 150):
#     tx_hash = contract_instance.functions.calculateAndStoreSquareSum(key_to_calculate).transact({
#         'from': '0xCeE4aecD7Af726D8338ef6e6eAdF12591E6e3b28',
#         'gas': 999999999999999  # You might need to adjust this value
#     })

#     # Wait for the transaction to be confirmed
#     tx_receipt = w3.eth.wait_for_transaction_receipt(tx_hash)
#     assert tx_receipt['status'] == 1

#     # Accumulate gas costs
#     total_gas_used += tx_receipt['gasUsed']

#     print(f"Data for key {key_to_calculate} successfully processed! Gas used: {tx_receipt['gasUsed']}")

# print(f"Total gas used for all transactions: {total_gas_used}")
