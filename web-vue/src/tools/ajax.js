var baseURL = 'http://192.168.1.6:8000'
// var baseURL = 'http://localhost:8000'
// 如果是生成环境，则无需设置baseURL
if (process.env.NODE_ENV === 'production') {
    baseURL = ''
}

function post(url, params, callback, err) {
    var xhr = new XMLHttpRequest();
    xhr.open('POST', baseURL + url, true);
    if (params) {
        var formData = new FormData()
        for (var key in params) {
            formData.append(key, params[key])
        }
        xhr.send(formData)
    } else {
        xhr.send()
    }
    xhr.onreadystatechange = function (data) {
        if (xhr.readyState === 4) {
            if ((xhr.status >= 200 && xhr.status < 300) || xhr.status === 304) {
                if (callback) {
                    callback(data.srcElement.response)
                }
            } else {
                if (err) {
                    err()
                }
            }
        }
    }
}



export { post, baseURL }