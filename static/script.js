
let selectedFileElement = null;
let selectedFile = null;
let currentPath = '';

const breadcrumbs = document.getElementById('breadcrumbs');
const fileList = document.getElementById('fileList');
const fileInput = document.getElementById('fileInput');

fileInput.addEventListener('change', (event) => {
    const files = event.target.files;
    for (let i = 0; i < files.length; i++) {
        uploadFile(files[i]);
    }
});

function uploadFile(file) {
    const xhr = new XMLHttpRequest();
    const formData = new FormData();
    formData.append('file', file);
    formData.append('path', currentPath);
    xhr.open('POST', '/api/upload', true);
    xhr.send(formData);
    xhr.onload = () => refreshFileList(currentPath);
}

document.body.addEventListener('dragover', (event) => {
    event.preventDefault();
    document.body.classList.add('dragging');
});

document.body.addEventListener('dragleave', (event) => {
    document.body.classList.remove('dragging');
});

document.body.addEventListener('drop', (event) => {
    event.preventDefault();
    document.body.classList.remove('dragging');
    const files = event.dataTransfer.files;
    for (let i = 0; i < files.length; i++) {
        uploadFile(files[i]);
    }
});

document.getElementById('downloadButton').addEventListener('click', async () => {
    if (selectedFile) {
        const response = await fetch('/api/download?name=' + encodeURIComponent(currentPath + '/' + selectedFile));
        const blob = await response.blob();
        const url = window.URL.createObjectURL(blob);
        const a = document.createElement('a');
        a.href = url;
        a.download = selectedFile;
        document.body.appendChild(a);
        a.click();
        a.remove();
    }
});

document.getElementById('addFolderButton').addEventListener('click', () => {
    const folderName = prompt('输入新的文件夹名称:');
    if (folderName) {
        const xhr = new XMLHttpRequest();
        xhr.open('POST', '/api/addFolder?name=' + encodeURIComponent(folderName) + '&path=' + encodeURIComponent(currentPath), true);
        xhr.onload = () => refreshFileList(currentPath);
        xhr.send();
    }
});

document.getElementById('refreshButton').addEventListener('click',
    () => refreshFileList(currentPath));

document.getElementById('deleteButton').addEventListener('click', () => {
    if (selectedFile) {
        if (confirm('确定要删除 ' + selectedFile + ' 吗？')) {
            const xhr = new XMLHttpRequest();
            xhr.open('POST', '/api/delete?name=' + encodeURIComponent(currentPath + '/' + selectedFile), true);
            xhr.onload = () => {
                if (xhr.status >= 200 && xhr.status < 300) {
                    refreshFileList(currentPath);
                    alert('删除成功');
                } else {
                    alert('删除失败: ' + xhr.responseText);
                }
            };
            xhr.onerror = () => {
                alert('删除失败: 网络错误');
            };
            xhr.send();
        }
    } else {
        alert('请选择要删除的文件');
    }
});

function refreshFileList(path) {
    currentPath = path; // 更新当前路径

    // 清除选中的文件
    selectedFileElement = null;
    selectedFile = null;

    // 更新面包屑
    breadcrumb.innerHTML = '';
    const pathParts = path.split('/');
    let pathAccumulator = '';
    for (let i = 0; i < pathParts.length; i++) {
        const pathPart = pathParts[i];
        if (pathPart) { // 跳过空字符串
            pathAccumulator += pathPart;
            const link = document.createElement('a');
            link.textContent = pathPart === '.' ? '根目录' : pathPart;;
            link.href = pathAccumulator;
            link.dataset.path = pathAccumulator;
            link.addEventListener('click', (event) => {
                event.preventDefault();
                refreshFileList(event.target.dataset.path);
            });
            breadcrumb.appendChild(link);
        }
        if (i < pathParts.length - 1) {
            breadcrumb.appendChild(document.createTextNode(' / '));
            pathAccumulator += '/';
        }
    }

    // 更新文件列表
    const xhr = new XMLHttpRequest();
    xhr.open('GET', '/api/files?path=' + encodeURIComponent(path), true);
    xhr.onload = function () {
        fileList.innerHTML = '';

        const files = JSON.parse(this.responseText);
        const table = document.createElement('table');
        const docFragment = document.createDocumentFragment(); // 创建文档片段

        for (let i = 0; i < files.length; i++) {
            const file = files[i];
            const row = document.createElement('tr');
            const nameCell = document.createElement('td');
            nameCell.textContent = file.name;
            const sizeCell = document.createElement('td');
            sizeCell.textContent = formatSize(file.size);
            const modTimeCell = document.createElement('td');
            modTimeCell.textContent = new Date(file.modTime).toLocaleString();
            row.appendChild(nameCell);
            row.appendChild(sizeCell);
            row.appendChild(modTimeCell);
            row.addEventListener('click', function () {
                if (selectedFileElement) {
                    selectedFileElement.classList.remove('selected');
                }
                this.classList.add('selected');
                selectedFile = file.name; // 更新为当前点击的文件名
                selectedFileElement = this; // 保存当前选中的元素
            });
            row.addEventListener('dblclick', function () {
                if (file.isDir) {
                    refreshFileList(currentPath + '/' + file.name);
                }
            });

            // 处理移动端的双击事件
            row.addEventListener('touchstart', function (e) {
                this.dataset.touchtime = Date.now();
            });
            row.addEventListener('touchend', function (e) {
                if (Date.now() - this.dataset.touchtime < 300) { // 300ms内的两次触摸被认为是双击
                    if (file.isDir) {
                        refreshFileList(currentPath + '/' + file.name);
                    }
                }
            });

            docFragment.appendChild(row); // 将行添加到文档片段中，而不是直接添加到表格中
        }
        table.appendChild(docFragment); // 一次性将所有行添加到表格中
        fileList.appendChild(table);
    };
    xhr.send();
}

function renameHandle(selectedFile) {
    const newName = prompt('输入新的文件名:', selectedFile);
    if (newName) {
        const xhr = new XMLHttpRequest();
        xhr.open('POST', '/api/rename?oldName=' + encodeURIComponent(currentPath + '/' + selectedFile) + '&newName=' + encodeURIComponent(newName), true);
        xhr.onload = () => refreshFileList(currentPath);
        xhr.send();
    }
}

// 新增：为重命名按钮添加点击事件监听器
document.getElementById('renameButton').addEventListener('click', () => {
    if (selectedFile) {
        renameHandle(selectedFile)
    }
});

document.addEventListener('keydown', function (event) {
    if (event.key == 'F2' && selectedFile) {
        renameHandle(selectedFile)
    }
});

function formatSize(size) {
    if (size === 0) {
        return 'NA';
    }
    const i = Math.floor(Math.log(size) / Math.log(1024));
    return (size / Math.pow(1024, i)).toFixed(2) * 1 + ' ' + ['B', 'kB', 'MB', 'GB', 'TB'][i];
}
refreshFileList('.');