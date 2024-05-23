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
                        <td>
                            <span v-if="item.isEditMode" style="display: flex; align-items: center;">
                                <input v-model="editFileName" type="text" class="editInput" ref="editInput"
                                    @keyup.enter="submitRename" @keyup.esc="cancelRename" @blur="cancelRename"
                                    @dblclick.stop>
                                <van-button icon="success" type="primary" size="mini" @click="submitRename" />
                                <van-button icon="cross" type="info" size="mini" @click="cancelRename" />
                            </span>
                            <span v-else>{{ item.name }}</span>
                        </td>
                        <td>{{ formatSize(item.size) }}</td>
                        <td>{{ item.modTime }}</td>
                    </tr>
                </tbody>
                <van-empty v-else description="文件列表为空" />
            </table>
        </div>

        <van-popup v-model="showPopup" position="top" :style="{ height: '100%' }" closeable>

            <div class="popup-title" v-if="selectItem">{{ selectItem.name }}</div>
            <pre class="popup-content" v-html="fileContent"></pre>
        </van-popup>
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
            showPopup: false,
            editFileName: '',
            fileContent: '',
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

        window.addEventListener('keydown', this.handleKeyDown);
    },
    beforeDestroy() {
        window.removeEventListener('keydown', this.handleKeyDown);
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
            var that = this

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
                } else if (/(.txt$|.md$|.log$)/i.test(lowercaseName)) {
                    post('/api/download?name=' + encodeURIComponent(this.currentPath + '/' + item.name), null, function (data) {
                        that.fileContent = data
                        that.showPopup = true
                    }, function () {
                        Toast('文件下载失败')
                    })
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
                    list[i].isEditMode = false
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
        handleKeyDown(event) {
            if (event.key === 'F2' && this.selectItem) {
                this.renameHandler();
                event.preventDefault(); // 阻止F2键的默认行为，如打开浏览器的查找功能
            }
        },
        submitRename() {
            if (this.editFileName.trim() !== '') {
                this.selectItem.isEditMode = false

                const selectedFile = this.selectItem.name
                const newName = this.editFileName.trim()
                if (newName) {
                    var that = this

                    var url = '/api/rename?oldName=' + encodeURIComponent(this.currentPath + '/' + selectedFile) + '&newName=' + encodeURIComponent(newName)
                    post(url, null, function () {
                        that.refreshHandler()
                    }, function () {
                        Toast('文件列表获取失败')
                    })
                }
            } else {
                this.selectItem.isEditMode = false
            }
            this.editFileName = ''
        },
        cancelRename() {
            this.selectItem.isEditMode = false
            this.editFileName = ''
        },
        renameHandler() {
            if (!this.selectItem) {
                Toast('请选择一个文件或文件夹')
                return
            }
            this.editFileName = this.selectItem.name
            this.selectItem.isEditMode = true

            var endIndex = this.selectItem.name.lastIndexOf('.')
            if (this.selectItem.isDir) {
                endIndex = this.selectItem.name.length
            }

            // 延迟执行，确保输入框被渲染
            this.$nextTick(() => {
                this.$refs.editInput[0].focus()
                this.$refs.editInput[0].setSelectionRange(0, endIndex)
            })
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

#breadcrumb a {
    font-size: 20px;
    line-height: 1.75;
}

.van-button:not(:last-child) {
    margin-right: 5px;
}

.editInput {
    height: 22px;
    margin-right: 5px;
}

.popup-title {
    height: 30px;
    width: 97.75%;
    border-bottom: 1px solid #dddddd;
    padding: 10px;
    font-size: 20px;
}

.popup-content {
    padding: 10px;
    width: 97%
}
</style>