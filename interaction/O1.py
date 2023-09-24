from web3 import Web3, AsyncWeb3

w3 = Web3(Web3.HTTPProvider('http://127.0.0.1:7545'))

address = '0xf566deEa66B871c83D0819E469219569362B1043'
abi = '''[
	{
		"inputs": [
			{
				"internalType": "int256",
				"name": "key",
				"type": "int256"
			},
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
	},
	{
		"inputs": [
			{
				"internalType": "int256[]",
				"name": "keys",
				"type": "int256[]"
			}
		],
		"name": "addArrays",
		"outputs": [
			{
				"internalType": "int256[]",
				"name": "resultGrad",
				"type": "int256[]"
			},
			{
				"internalType": "int256[]",
				"name": "resultBeta",
				"type": "int256[]"
			},
			{
				"internalType": "int256[]",
				"name": "resultSigma",
				"type": "int256[]"
			}
		],
		"stateMutability": "view",
		"type": "function"
	},
	{
		"inputs": [
			{
				"internalType": "int256",
				"name": "key",
				"type": "int256"
			}
		],
		"name": "getBeta",
		"outputs": [
			{
				"internalType": "int256[]",
				"name": "",
				"type": "int256[]"
			}
		],
		"stateMutability": "view",
		"type": "function"
	},
	{
		"inputs": [
			{
				"internalType": "int256",
				"name": "key",
				"type": "int256"
			}
		],
		"name": "getGrad",
		"outputs": [
			{
				"internalType": "int256[]",
				"name": "",
				"type": "int256[]"
			}
		],
		"stateMutability": "view",
		"type": "function"
	},
	{
		"inputs": [
			{
				"internalType": "int256",
				"name": "key",
				"type": "int256"
			}
		],
		"name": "getSigma",
		"outputs": [
			{
				"internalType": "int256[]",
				"name": "",
				"type": "int256[]"
			}
		],
		"stateMutability": "view",
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




# key = 7450
# tx_hash = contract_instance.functions.storeData(key, grad_list, beta_list, sigma_list).transact({'from': '0x5acc3435F5436aE463cd60E456E76435A336aedB'})

# # Wait for the transaction to be mined
# tx_receipt = w3.eth.wait_for_transaction_receipt(tx_hash)

# # Check if the transaction was successful
# assert tx_receipt['status'] == 1

# print(f"Data successfully stored for key {key}!")












# ----------------------------------------------------------------------------------------------------------------------#

# ----------------------------------------------------------------------------------------------------------------------#


# def validate_int256_numbers(numbers):
#     for num in numbers:
#         assert isinstance(num, int), f"Number {num} is not an integer!"
#         assert -2**255 <= num < 2**255, f"Number {num} is out of range for int256!"

# def transact_data(key, start_idx, end_idx, grad, beta, sigma):
#     # Validate the numbers before transacting
#     validate_int256_numbers(grad)
#     validate_int256_numbers(beta)
#     validate_int256_numbers(sigma)

#     print(f"Storing array indices: {start_idx} to {end_idx-1}")
#     tx_hash = contract_instance.functions.storeData(key, grad, beta, sigma).transact({'from': '0xCeE4aecD7Af726D8338ef6e6eAdF12591E6e3b28','gas':999999999999999})

#     # Wait for the transaction to be mined
#     tx_receipt = w3.eth.wait_for_transaction_receipt(tx_hash)

#     # Check if the transaction was successful
#     assert tx_receipt['status'] == 1

#     gas_used = tx_receipt['gasUsed']
#     print(f"Data successfully stored for key {key}! Gas used: {gas_used}")
#     return gas_used
    

# # Starting key
# key = 100

# # Initialize a variable to store total gas used
# total_gas_used = 0

# # Process each list in batches of 
# for i in range(0, len(grad_list), 50):
#     grad_batch = grad_list[i:i+50]
#     beta_batch = beta_list[i:i+50]
#     sigma_batch = sigma_list[i:i+50]

#     gas_used_for_batch = transact_data(key, i, i+len(grad_batch), grad_batch, beta_batch, sigma_batch)
#     total_gas_used += gas_used_for_batch
#     key += 1  # Increment the key for the next batch

# print(f"Total gas used for all transactions: {total_gas_used}")






# ----------------------------------------------------------------------------------------------------------------------#

# ----------------------------------------------------------------------------------------------------------------------#

# Assuming you've already defined the address variable earlier in your code
sender_address = "0xCeE4aecD7Af726D8338ef6e6eAdF12591E6e3b28"

# Define the arrays you want to add
arrays_to_add = [[i,i,i,i] for i in range(100, 248)]

# Initialize a variable to store total gas used
total_gas_used = 0

for array in arrays_to_add:
    # Call the addArrays function to transact
    tx_hash = contract_instance.functions.addArrays(array).transact({'from': sender_address})
    # Wait for the transaction to be mined
    tx_receipt = w3.eth.wait_for_transaction_receipt(tx_hash)
    
    # Check if the transaction was successful
    if tx_receipt['status'] == 1:
        gas_used = tx_receipt['gasUsed']
        total_gas_used += gas_used
        print(f"Array {array} added successfully! Gas used: {gas_used}")
    else:
        print(f"Failed to add array {array}!")

print(f"Total gas used for all transactions: {total_gas_used}")
















