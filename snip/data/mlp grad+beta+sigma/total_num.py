import re

# 定义一个正则表达式来匹配这种格式的数字
pattern = r"\d+\.\d+e[+-]\d+"

def count_numbers_in_file(filename):
    with open(filename, 'r') as f:
        content = f.read()
        # 使用正则表达式查找所有匹配的数字
        numbers = re.findall(pattern, content)
        return len(numbers)

filename = "betafc1.txt"
count = count_numbers_in_file(filename)
print(f"文件中有 {count} 个这种格式的数字。")
