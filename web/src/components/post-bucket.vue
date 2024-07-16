<template>
  <div class="link-wrap">
    <div class="link-item" v-if="postContentText.bucket.length > 0">
      <n-icon class="hash-link"><archive-outline /></n-icon>
      <div class="link-txt-wrap">
        <span @click="openEditModel()" class="link-txt">{{ postContentText.bucket + "/" +
          postContentText.object_dir}}</span>
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
        <n-select v-model:value="postContentText.bucket" filterable placeholder="选一个bucket" :options="options" />
      </n-input-group>

      <n-input-group style="margin-bottom: 2px;">
        <n-button :style="{ width: '12%' }" disabled>
          dir
        </n-button>
        <n-input type="text" v-model:value="postContentText.object_dir" />
      </n-input-group>

      <n-input-group style="margin-bottom: 2px;">
        <n-button :style="{ width: '12%' }" disabled>
          post_id
        </n-button>
        <n-input-number disabled type="text" v-model:value="postContentText.post_id" />
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


import {
  editPostBucket
} from '@/api/post';


const props = withDefaults(defineProps<{
    items: Item.PostItemProps[]
}>(), {
  items: () => []
});


// cyw 内容编辑
const postContentText = ref({
  "post_id": 0,
  "bucket": "",
  "object_dir": "",
  // "key": "",
})

for (let i = 0; i < props.items.length; i++) {
  const element = props.items[i];
  if (element.content.indexOf("http") >= 0) {
    var a = new URL(element.content)
    var temps = a.pathname.split("/")

    postContentText.value.post_id = element.post_id
    // bucketInfo.key = temps[temps.length-1]

    temps = temps.slice(1);
    temps = temps.slice(0, -1);
    postContentText.value.bucket = temps[0]
    temps = temps.slice(1);
    postContentText.value.object_dir = temps.join("/")

    console.log(postContentText.value)
    break;
  }

}

// 编辑
const showEditModal = ref(false);

const execEditAction = () => {
  editPostBucket(postContentText.value)
    .then((_res) => {
      window.$message.success('修改成功');
      showEditModal.value = false;


      postContentText.value.post_id = _res.post_id
      postContentText.value.bucket = _res.bucket
      postContentText.value.object_dir = _res.object_dir
    })
    .catch((_err) => {
      loading.value = false;
    });
};

const openEditModel = () =>{
  showEditModal.value = true;
}

// buckets
const  options= [
      {
        label: 'social-station',
        value: 'social-station'
      },
      {
        label: 'photo-station',
        value: 'photo-station'
      },
      {
        label: 'topic-station',
        value: 'topic-station'
      },
      {
        label: 'paopao-station',
        value: 'paopao-station'
      },
      {
        label: 'video-station',
        value: 'video-station'
      },
]


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