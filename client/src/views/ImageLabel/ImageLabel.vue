<template>
    <div class="container" id="container">
        <div class="left">
            <form action="">
                <label for="images">图片:</label>
                <div style="display: flex;justify-content: space-between;">
                    <a-button @click="selectImages">选择图片</a-button>
                    <input type="file" id="images" name="images[]" accept="image/jpeg, image/png" multiple
                           style="display: none"/>
                    <a-button name="cropImages" id="cropImages">裁剪并保存</a-button>
                </div>
                <a-divider>图片列表</a-divider>
                <label for="imageSearch">搜索:</label>
                <div style="margin-bottom: 10px">
                    <a-input type="text" id="imageSearch" name="imageSearch"/>
                </div>
                <label for="imageList"></label>
                <select name="imageList" id="imageList" size="10" multiple></select>

                <div id="imageInformation"></div>
                <a-divider>类别列表</a-divider>

                <label for="classes">类别:</label>
                <div style="margin-bottom: 10px">
                    <a-button @click="selectFile">选择类别文件</a-button>
                    <input type="file" id="classes" name="classes" accept="text/plain" style="display: none"/>
                </div>

                <label for="classList"></label>
                <select name="classList" id="classList" size="10" multiple style="width: 100%"></select>

                <div id="bboxInformation"></div>
                <a-divider>标注文件相关</a-divider>

                <label for="bboxes">Bboxes:</label>
                <div style="display: flex;justify-content: space-between;margin-bottom: 10px">
                    <a-button @click="selectAnnotationFile">选择标注文件</a-button>
                    <input type="file" id="bboxes" name="bboxes[]" accept="text/plain, application/zip" disabled
                           multiple style="display: none"/>
                    <a-button name="restoreBboxes" id="restoreBboxes" disabled>重置</a-button>
                </div>

                <a-divider style="margin: 8px 0"></a-divider>
                <a-button name="saveBboxes" id="saveBboxes">存储为YOLO文件</a-button>
                <a-divider style="margin: 8px 0"></a-divider>
                <a-button name="saveCocoBboxes" id="saveCocoBboxes">存储为COCO文件</a-button>
                <a-divider style="margin: 8px 0"></a-divider>
                <div id="voc" style="margin-bottom: 10px">
                    <label for="vocFolder">Voc 文件夹:</label>
                    <div style="margin-bottom: 10px">
                        <a-input type="text" id="vocFolder" name="vocFolder" value="data" style="margin-bottom: 10px"/>
                        <a-button name="saveVocBboxes" id="saveVocBboxes">存储为VOC文件</a-button>
                    </div>
                </div>
                <div id="coco"></div>
                <a-divider>操作提示</a-divider>

                <div id="description">
                    <ul>
                        <li>Mouse WHEEL - zoom in/out image</li>
                        <li>Mouse RIGHT BUTTON - pan image</li>
                        <li>Arrows LEFT and RIGHT - cycle images</li>
                        <li>Arrows Up and DOWN - cycle classes</li>
                        <li>Key DELETE - remove selected Bbox</li>
                    </ul>
                </div>
            </form>
        </div>

        <div class="right" id="right">
            <canvas id="canvas"></canvas>
        </div>
    </div>
</template>

<script lang="ts">
    import {Component, Vue} from 'vue-property-decorator'
    import './boobs/ybat'

    @Component
    export default class ImageLabel extends Vue {
        private created(): void {

        }

        private mounted(): void {
            window.dispatchEvent(new Event('image-label-mounted'))
        }

        private selectImages(): void {
            (document.getElementById('images') as HTMLInputElement).click()
        }

        private selectFile(): void {
            (document.getElementById('classes') as HTMLInputElement).click()
        }

        private selectAnnotationFile(): void {
            (document.getElementById('bboxes') as HTMLInputElement).click()
        }
    }
</script>

<style scoped>
    label {
        font-weight: bold;
    }

    form {
        height: 100%
    }

    hr {
        width: 100%;
    }

    .container {
        min-height: 100%;
        width: 100%;
        display: flex;
        flex-direction: row;
    }

    .left {
        min-width: 300px;
        width: 30%;
        height: 100%;
        margin-right: 10px;
    }

    .right {
        width: 100%;
        margin-left: 10px;
        border-left: 1px solid #cccccc;
    }

    #imageList, #classList {
        margin-bottom: 10px;
        width: 100%;
        border-radius: 4px;
        border: 1px solid #d9d9d9;
        transition: all 0.3s;
    }

    #imageList:focus, #classList:focus {
        border-color: #40a9ff;
        outline: 0;
        box-shadow: 0 0 0 2px rgba(24, 144, 255, 0.2);
    }

    >>> .ant-divider-with-text::before, .ant-divider-with-text::after {
        border-top: 1px solid #cccccc !important;
    }
</style>
