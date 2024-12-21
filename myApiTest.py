import uuid

from flask import Flask, request, jsonify

app = Flask(__name__)


@app.route('/task/list', methods=['GET'])
def task_list():
    print(request.headers)
    # 检查Cookie是否正确
    heat_cookie = request.cookies.get('heat')
    if heat_cookie != '123456':
        return jsonify({"error": "Missing or invalid cookie"}), 400

    # 检查URL参数taskid是否存在且等于43234344
    taskid = request.args.get('taskid')
    if not taskid or taskid != '43234344':
        return jsonify({"error": "Missing or invalid taskid"}), 400

    # 如果所有验证都通过，则返回成功响应
    response = jsonify({"message": "Request is valid", "taskid": taskid})

    # 增加一个响应头 traceId
    trace_id = str(uuid.uuid4())  # 生成一个唯一的traceId
    response.headers['traceId'] = trace_id

    return response, 200


@app.route('/task/detail', methods=['POST'])
def task_detail():
    # 检查请求头中的role-id是否为admin
    print(request.headers)
    role_id = request.headers.get('role-id')
    if role_id != 'admin':
        return jsonify({"error": "Missing or invalid role_id"}), 400

    # 检查请求头中的group-id是否为123
    group_id = request.headers.get('group-id')
    if group_id != '123':
        return jsonify({"error": "Missing or invalid group_id"}), 400

    # 检查Cookie是否正确
    heat_cookie = request.cookies.get('heat')
    if heat_cookie != '123456':
        return jsonify({"error": "Missing or invalid cookie"}), 400

    # 检查表单参数taskid是否存在且等于43234344
    taskid = request.form.get('taskid')
    if not taskid or taskid != '43234344':
        return jsonify({"error": "Missing or invalid taskid"}), 400

    # 如果所有验证都通过，则返回成功响应
    response = jsonify({"message": "Request is valid", "taskid": taskid})

    # 增加一个响应头 traceId
    trace_id = str(uuid.uuid4())  # 生成一个唯一的traceId
    response.headers['traceId'] = trace_id

    return response, 200


@app.route('/task/order', methods=['POST'])
def task_order():
    # 检查请求头中的role-id是否为admin
    print(request.headers)
    role_id = request.headers.get('role-id')
    if role_id != 'admin':
        return jsonify({"error": "Missing or invalid role_id"}), 400

    # 检查请求头中的group-id是否为123
    group_id = request.headers.get('group-id')
    if group_id != '123':
        return jsonify({"error": "Missing or invalid group_id"}), 400

    # 检查Cookie是否正确
    heat_cookie = request.cookies.get('heat')
    if heat_cookie != '123456':
        return jsonify({"error": "Missing or invalid cookie"}), 400

    try:
        # 获取并解析JSON请求体
        data = request.get_json()

        # 检查JSON参数aa和bb是否存在且等于预期值
        if not data or 'aa' not in data or 'bb' not in data:
            return jsonify({"error": "Missing JSON parameters"}), 400

        aa = data.get('aa')
        bb = data.get('bb')

        if aa != 123 or bb != 456:
            return jsonify({"error": "Invalid JSON parameter values"}), 400

    except Exception as e:
        return jsonify({"error": "Invalid JSON format", "details": str(e)}), 400

    # 如果所有验证都通过，则返回成功响应
    response = jsonify({"message": "Request is valid", "data": {"aa": aa, "bb": bb}})

    # 增加一个响应头 traceId
    trace_id = str(uuid.uuid4())  # 生成一个唯一的traceId
    response.headers['traceId'] = trace_id

    return response, 200


if __name__ == '__main__':
    app.run(debug=True)