<template>
  <div class="link-wrap">
    <div class="link-item" v-if="bucketInfo.bucket.length > 0">
      <n-icon class="hash-link"><archive-outline /></n-icon>
      <div class="link-txt-wrap">
        <span @click="openEditModel()" class="link-txt">{{ bucketInfo.bucket + "/" + bucketInfo.dir}}</span>
      </div>
    </div>

    <!-- 编辑content -->
    <n-modal v-model:show="showEditModal" :mask-closable="false" class="custom-card" style="width:600px;" preset="card"
      title="修改" positive-text="确认" negative-text="取消" @positive-click="execEditAction">
      <template #header-extra>
      </template>


      <n-input-group style="margin-bottom: 2px;">
        <n-button :style="{ width: '12%'}" disabled>
          bucket
        </n-button>
        <n-input type="text" v-model:value="postContentText.bucket" />
      </n-input-group>

      <n-input-group style="margin-bottom: 2px;">
        <n-button :style="{ width: '12%' }" disabled>
          dir
        </n-button>
        <n-input type="text" v-model:value="postContentText.dir" />
      </n-input-group>

      <n-input-group style="margin-bottom: 2px;">
        <n-button :style="{ width: '12%' }" disabled>
          post_id
        </n-button>
        <n-input disabled type="text" v-model:value="postContentText.post_id" />
      </n-input-group>


      <template #footer>
        <n-button @click="execEditAction" type="primary" style="float:right;">
          保存
        </n-button>
      </template>
    </n-modal>


  </div>
</template>

<script setup lang="ts">
import { h, ref, onMounted, computed } from 'vue';
import { ArchiveOutline } from '@vicons/ionicons5';
const props = withDefaults(defineProps<{
    items: Item.PostItemProps[]
}>(), {
  items: () => []
});

var bucketInfo = {
  "post_id": 0,
  "bucket": "",
  "dir": "",
  "key": "",
}

for (let i = 0; i < props.items.length; i++) {
  const element = props.items[i];
  if (element.content.indexOf("http") >= 0) {
    var a = new URL(element.content)
    var temps = a.pathname.split("/")

    bucketInfo.post_id = element.post_id
    bucketInfo.key = temps[temps.length-1]

    temps = temps.slice(1);
    temps = temps.slice(0, -1);
    bucketInfo.bucket = temps[0]
    temps = temps.slice(1);
    bucketInfo.dir = temps.join("/")
    break;
  }
}
console.log(bucketInfo)


// 编辑
const showEditModal = ref(false);

// cyw 内容编辑
const postContentText = ref({
  "post_id": 0,
  "bucket": "",
  "dir": "",
  "key": "",
})

const execEditAction = () => {
  window.$message.success('修改成功');
  // editPostText(postContentText.value)
  //   .then((_res) => {
  //     showEditModal.value = false;
  //   })
  //   .catch((_err) => {
  //     loading.value = false;
  //   });
};

const openEditModel = () =>{
  showEditModal.value = true;
  postContentText.value = bucketInfo
}

</script>

<style lang="less" scoped>
.link-wrap {
    margin-bottom: 10px;
    position: relative;
    .link-item {
        height: 22px;
        display: flex;
        align-items: center;
        position: relative;
        .link-txt-wrap {
            left: calc(1em + 4px);
            width: calc(100% - 1em);
            overflow: hidden;
            white-space: nowrap;
            text-overflow: ellipsis;
            position: absolute;
            .hash-link {
                .link-txt {
                    word-break: break-all;
                }
            }
        }
    }
}
</style>