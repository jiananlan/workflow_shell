import requests
import base64
import json

# GitHub仓库信息
username = 'jiananlan'  # GitHub用户名
repo = 'test'  # 仓库名
file_path = 'test.txt'  # 文件路径
branch = 'main'  # 分支名
token = '！！！！填这里！！！！'  # 你的GitHub个人访问令牌

# API的基础URL
api_url = f'https://api.github.com/repos/{username}/{repo}/contents/{file_path}?ref={branch}'


# 获取文件内容
def get_file_content():
    headers = {'Authorization': f'token {token}'}
    response = requests.get(api_url, headers=headers)

    if response.status_code == 200:
        file_info = response.json()
        content = base64.b64decode(file_info['content']).decode('utf-8')
        return content
    else:
        print(f"Error fetching file: {response.status_code}")
        return None


# 更新文件内容
def update_file_content(new_content):
    # 获取当前文件信息（如sha值）
    headers = {'Authorization': f'token {token}'}
    response = requests.get(api_url, headers=headers)

    if response.status_code == 200:
        file_info = response.json()
        sha = file_info['sha']  # 获取当前文件的sha值

        # 准备新的文件内容
        new_content_base64 = base64.b64encode(new_content.encode('utf-8')).decode('utf-8')

        # 请求更新文件
        update_url = f'https://api.github.com/repos/{username}/{repo}/contents/{file_path}'
        data = {
            'message': 'Update text_content.txt',  # 提交信息
            'content': new_content_base64,  # 新的文件内容（Base64编码）
            'sha': sha,  # 当前文件的sha值
            'branch': branch  # 目标分支
        }

        response = requests.put(update_url, headers=headers, json=data)

        if response.status_code == 200:
            print("File updated successfully!")
        else:
            print(f"Error updating file: {response.status_code}")
    else:
        print(f"Error fetching file: {response.status_code}")


# 示例使用
if __name__ == "__main__":
    # 获取文件内容
    current_content = get_file_content()
    if current_content:
        print("Current content:",end='')
        print(current_content)

    # 更新文件内容
    new_content = "This is the new content for the API."
    while True:
        update_file_content(input())

