<template>
    <div ref="main"></div>
</template>

<script>
import * as monaco from "monaco-editor";

export default {
    name: "CodeEditor",
    data() {
        return {
            monacoEditor: null,
        };
    },
    props: {
        language: String,
        value: String
    },
    created() {

    },
    watch: {
        value(newCode) {
            this.monacoEditor.setValue(newCode);
        }
    },
    mounted() {
        this.$nextTick(() => {
            this.initMonaco();
        })
    },
    methods: {
        initMonaco() {
            this.monacoEditor = monaco.editor.create(this.$refs.main, {
                theme: "vs-dark", // 主题
                value: this.value, // 默认显示的值
                language: this.language,
                folding: true, // 是否折叠
                foldingHighlight: true, // 折叠等高线
                foldingStrategy: "auto", // 折叠方式
                showFoldingControls: "always", // 是否一直显示折叠
                disableLayerHinting: true, // 等宽优化
                emptySelectionClipboard: false, // 空选择剪切板
                selectionClipboard: false, // 选择剪切板
                automaticLayout: true, // 自动布局
                codeLens: true, // 代码镜头
                scrollBeyondLastLine: false, // 滚动完最后一行后再滚动一屏幕
                colorDecorators: true, // 颜色装饰器
                accessibilitySupport: "on", // 辅助功能支持"auto" | "off" | "on"
                lineNumbers: "on", // 行号 取值： "on" | "off" | "relative" | "interval" | function
                lineNumbersMinChars: 4, // 行号最小字符   number
                enableSplitViewResizing: false,
                readOnly: false, //是否只读  取值 true | false
                fontSize: 18
            });
            this.monacoEditor.onDidChangeModelContent(() => {
                this.$emit('change', this.monacoEditor.getValue());
            });
        }

    },
};
</script>