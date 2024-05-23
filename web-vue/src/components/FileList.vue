<template>
    <div class="hello">

        <h2>局域网文件共享 Shared2</h2>
        <input type="file" @change="uploadsHandler" id="fileInput" multiple hidden>
        <div id="toolbar">
            <van-button size="small" type="info" @click="downloadHandler">&#x1F4E5; 下载文件</van-button>
            <van-button size="small" type="info" onclick="fileInput.click()">&#x1F4E4; 上传文件</van-button>
            <van-button size="small" type="info" @click="refreshHandler">&#x21BB; 刷新</van-button>
            <van-button size="small" type="info" @click="addFolderHandler">&#x1F4C1; 添加文件夹</van-button>
            <van-button size="small" type="info" @click="renameHandler">&#x270E; 重命名</van-button>
            <van-button size="small" type="info" @click="deleteHandler">&#x1F5D1; 删除</van-button>
            <van-button size="small" type="info" @click="intoHandler(selectItem)">&#x1F4CD; 打开</van-button>
        </div>

        <div id="breadcrumb">
            <span v-for="(item, index) in breadcrumbs" :key="index">
                <a v-if="item !== '.'" href="javascript:void(0)" @click="gotoHandler(index)">{{ item }}</a>
                <a v-else href="javascript:void(0)" @click="gotoHandler(0)">/所有文件</a>

                <span v-if="index < breadcrumbs.length - 1"> / </span>
            </span>
        </div>

        <div id="fileList">
            <table>
                <thead>
                    <tr>
                        <th>文件名</th>
                        <th>大小</th>
                        <th>修改时间</th>
                    </tr>
                </thead>
                <tbody v-if="files.length > 0">
                    <tr v-for="item in files" v-bind:key="item.name" @dblclick="intoHandler(item)"
                        @click="selectHandler(item)" v-bind:class="{ 'selected': item.isActive }">
                        <td>{{ item.name }}</td>
                        <td>{{ formatSize(item.size) }}</td>
                        <td>{{ item.modTime }}</td>
                    </tr>
                </tbody>
                <van-empty v-else description="文件列表为空" />
            </table>
        </div>
    </div>
</template>

<script>
import { ImagePreview, Toast } from 'vant';
import 'vant/lib/index.css';
import { post, baseURL } from '../tools/ajax'
export default {
    name: 'HelloWorld',
    components: {
        [ImagePreview.Component.name]: ImagePreview.Component,
    },
    props: {
        msg: String
    },
    data: function () {
        return {
            currentPath: '.',
            breadcrumbs: ['.'],
            selectItem: null,
            files: []
        }
    },
    mounted() {
        var that = this

        this.refreshHandler()

        // 添加拖拽事件
        document.body.addEventListener('dragover', function (event) {
            event.preventDefault();
            document.body.classList.add('dragging');
        });
        document.body.addEventListener('dragleave', function () {
            document.body.classList.remove('dragging');
        });
        document.body.addEventListener('drop', function (event) {
            event.preventDefault();
            document.body.classList.remove('dragging');
            const files = event.dataTransfer.files;
            for (let i = 0; i < files.length; i++) {
                that.uploadHandler(files[i]);
            }
        });
    },
    methods: {
        selectHandler(item) {
            if (this.selectItem) {
                this.selectItem.isActive = false
            }

            this.selectItem = item
            this.selectItem.isActive = true
        },
        gotoHandler(index) {
            if (index < this.breadcrumbs.length - 1) {
                this.breadcrumbs = this.breadcrumbs.slice(0, index + 1)
                this.currentPath = this.breadcrumbs.join('/')
                this.refreshHandler()
            }
        },
        intoHandler(item) {
            if (item == null) {
                Toast('请选择文件或文件夹')
                return
            }

            if (item.isDir) {
                var path = this.currentPath + '/' + item.name
                this.currentPath = path
                this.breadcrumbs = path.split('/')
                this.refreshHandler()
            } else {
                const lowercaseName = item.name.toLowerCase();
                if (/(.jpg$|.png$|.jpeg$|.gif$)/i.test(lowercaseName)) {
                    var url = baseURL + '/api/download?name=' + encodeURIComponent(this.currentPath + '/' + item.name)
                    this.imgPreview(url)
                    return
                } else {
                    Toast('暂不支持预览该文件')
                }
            }
        },
        refreshHandler() {
            this.selectItem = null

            var that = this
            post('/api/files?path=' + encodeURIComponent(this.currentPath), null, function (data) {
                var list = JSON.parse(data)
                for (var i = 0; i < list.length; i++) {
                    list[i].isActive = false
                    if (list[i].modTime) {
                        list[i].modTime = new Date(list[i].modTime).toLocaleString();
                    }
                }
                that.files = list
            }, function () {
                Toast('文件列表获取失败')
            })
        },
        deleteHandler() {
            if (!this.selectItem) {
                Toast('请选择一个文件或文件夹')
                return
            }

            var selectedFile = this.selectItem.name
            if (confirm('确定要删除 ' + selectedFile + ' 吗？')) {
                var that = this

                var url = '/api/delete?name=' + encodeURIComponent(this.currentPath + '/' + selectedFile)
                post(url, null, function () {
                    that.refreshHandler()
                }, function () {
                    Toast('文件列表获取失败')
                })
            }
        },
        addFolderHandler() {
            const folderName = prompt('输入新的文件夹名称:');
            if (folderName) {
                var that = this

                var url = '/api/addFolder?name=' + encodeURIComponent(folderName) + '&path=' + encodeURIComponent(this.currentPath)
                post(url, null, function () {
                    that.refreshHandler()
                }, function () {
                    Toast('文件列表获取失败')
                })
            }
        },
        renameHandler() {
            if (!this.selectItem) {
                Toast('请选择一个文件或文件夹')
                return
            }

            const selectedFile = this.selectItem.name
            const newName = prompt('输入新的文件名:', selectedFile);
            if (newName) {
                var that = this

                var url = '/api/rename?oldName=' + encodeURIComponent(this.currentPath + '/' + selectedFile) + '&newName=' + encodeURIComponent(newName)
                post(url, null, function () {
                    that.refreshHandler()
                }, function () {
                    Toast('文件列表获取失败')
                })
            }
        },
        downloadHandler() {
            if (!this.selectItem) {
                Toast('请选择一个文件或文件夹')
                return
            }
            var url = baseURL + '/api/download?name=' + encodeURIComponent(this.currentPath + '/' + this.selectItem.name)
            window.open(url, '_blank',)
        },
        uploadHandler(file) {
            var that = this
            var url = '/api/upload'
            var data = {
                'file': file,
                'path': this.currentPath
            }
            post(url, data, function () {
                that.refreshHandler()
            }, function () {
                Toast('文件列表获取失败')
            })
        },
        uploadsHandler() {
            var files = document.getElementById('fileInput').files
            for (var i = 0; i < files.length; i++) {
                this.uploadHandler(files[i]);
            }
        },
        imgPreview(url) {
            ImagePreview({
                images: [url],
                closeable: true,
            });
        },
        formatSize(size) {
            if (size === 0) {
                return 'NA';
            }
            const i = Math.floor(Math.log(size) / Math.log(1024));
            return (size / Math.pow(1024, i)).toFixed(2) * 1 + ' ' + ['B', 'kB', 'MB', 'GB', 'TB'][i];
        }
    }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
html,
body {
    height: calc(100% - 16px);
}

table {
    border-collapse: collapse;
    width: 100%;
}

td,
th {
    border: 1px solid #dddddd;
    text-align: left;
    padding: 8px;
}

tbody tr:hover {
    background-color: #f5f5f5;
}

tbody tr:nth-child(even) {
    background-color: #f2f2f2;
}

tbody tr.selected {
    outline: 2px dashed green;
}

.dragging {
    background-color: #eee;
    border: 3px dashed #aaa;
}

#breadcrumb a {
    font-size: 20px;
    line-height: 1.75;
}

.van-button:not(:last-child) {
    margin-right: 5px;
}
</style>