<template>
  <div class="post-item">
    <n-thing content-indented>
      <template #avatar>
        <n-avatar round :size="30" :src="post.user.avatar" />
      </template>
      <template #header>
        <span class="nickname-wrap">
          <router-link @click.stop class="username-link" :to="{
            name: 'user',
            query: { s: post.user.username },
          }">
            {{ post.user.nickname }}
          </router-link>
        </span>
        <span class="username-wrap"> @{{ post.user.username }} </span>
        <n-tag v-if="post.is_top" class="top-tag" type="warning" size="small" round>
          置顶
        </n-tag>
        <n-tag v-if="post.visibility == 1" class="top-tag" type="error" size="small" round>
          私密
        </n-tag>
        <n-tag v-if="post.visibility == 2" class="top-tag" type="info" size="small" round>
          好友可见
        </n-tag>
      </template>
      <template #header-extra>

        <div class="item-header-extra">
          <span class="timestamp">
            {{ post.ip_loc ? post.ip_loc + ' · ' : post.ip_loc }}
            {{ formatPrettyDate(post.created_on) }}
          </span>



          <div class="options" v-if="store.state.userInfo.is_admin ||
            store.state.userInfo.id === post.user.id
            ">



            <n-dropdown placement="bottom-end" trigger="click" size="small" :options="adminOptions"
              @select="handlePostAction">
              <n-button quaternary circle>
                <template #icon>
                  <n-icon>
                    <more-horiz-filled />
                  </n-icon>
                </template>
              </n-button>
            </n-dropdown>
          </div>
        </div>

        <!-- 删除确认 -->
        <n-modal v-model:show="showDelModal" :mask-closable="false" preset="dialog" title="提示" content="确定删除该泡泡动态吗？"
          positive-text="确认" negative-text="取消" @positive-click="execDelAction" />
        <!-- 锁定确认 -->
        <n-modal v-model:show="showLockModal" :mask-closable="false" preset="dialog" title="提示" :content="'确定' +
          (post.is_lock ? '解锁' : '锁定') +
          '该泡泡动态吗？'
          " positive-text="确认" negative-text="取消" @positive-click="execLockAction" />
        <!-- 置顶确认 -->
        <n-modal v-model:show="showStickModal" :mask-closable="false" preset="dialog" title="提示" :content="'确定' +
          (post.is_top ? '取消置顶' : '置顶') +
          '该泡泡动态吗？'
          " positive-text="确认" negative-text="取消" @positive-click="execStickAction" />
        <!-- 修改可见度确认 -->
        <n-modal v-model:show="showVisibilityModal" :mask-closable="false" preset="dialog" title="提示" :content="'确定将该泡泡动态可见度修改为' +
          (tempVisibility == 0 ? '公开' : (tempVisibility == 1 ? '私密' : '好友可见')) +
          '吗？'
          " positive-text="确认" negative-text="取消" @positive-click="execVisibilityAction" />




      </template>
      <template #description v-if="post.texts.length > 0">
        <span v-for="content in post.texts" :key="content.id" class="post-text" @click.stop="doClickText($event, post.id)"
          v-html="parsePostTag(content.content).content"></span>
      </template>

      <template #footer>
        <post-attachment v-if="post.attachments.length > 0" :attachments="post.attachments" />
        <post-attachment v-if="post.charge_attachments.length > 0" :attachments="post.charge_attachments"
          :price="post.attachment_price" />
        <post-image v-if="post.imgs.length > 0" :imgs="post.imgs" />
        <post-video v-if="post.videos.length > 0" :videos="post.videos" />
        <post-link v-if="post.links.length > 0" :links="post.links" />
      </template>
      <template #action>
        <n-space justify="space-between">
          <div class="opt-item">
            <n-icon size="18" class="opt-item-icon">
              <heart-outline />
            </n-icon>
            {{ post.upvote_count }}
          </div>
          <div class="opt-item" @click.stop="goPostDetail(post.id)">
            <n-icon size="18" class="opt-item-icon">
              <chatbox-outline />
            </n-icon>
            {{ post.comment_count }}
          </div>
          <div class="opt-item">
            <n-icon size="18" class="opt-item-icon">
              <bookmark-outline />
            </n-icon>
            {{ post.collection_count }}
          </div>
        </n-space>
      </template>
    </n-thing>
  </div>
</template>

<script setup lang="ts">
import { h, ref, onMounted, computed } from 'vue';
import type { Component } from 'vue'
import { NIcon } from 'naive-ui'
import { useStore } from 'vuex';
import { useRouter } from 'vue-router';
import { formatPrettyDate } from '@/utils/formatTime';
import { parsePostTag } from '@/utils/content';
import type { DropdownOption } from 'naive-ui';
import { VisibilityEnum } from '@/utils/IEnum';

import {
  Heart,
  Bookmark,
  ShareSocialOutline,
  PushOutline,
  TrashOutline,
  LockClosedOutline,
  LockOpenOutline,
  EyeOutline,
  EyeOffOutline,
  PersonOutline,
  FlameOutline,
  HeartOutline,
  BookmarkOutline,
  ChatboxOutline,
} from '@vicons/ionicons5';

import { MoreHorizFilled } from '@vicons/material';

import {
  getPostStar,
  postStar,
  getPostCollection,
  postCollection,
  deletePost,
  lockPost,
  stickPost,
  highlightPost,
  visibilityPost
} from '@/api/post';


const showDelModal = ref(false);
const showLockModal = ref(false);
const showStickModal = ref(false);
const showHighlightModal = ref(false);
const showVisibilityModal = ref(false);
const loading = ref(false);
const tempVisibility = ref<VisibilityEnum>(VisibilityEnum.PUBLIC);

const emit = defineEmits<{
  (e: 'reload'): void;
}>();



const router = useRouter();
const store = useStore();
const props = withDefaults(defineProps<{
  post: Item.PostProps,
}>(), {});

const renderIcon = (icon: Component) => {
  return () => {
    return h(NIcon, null, {
      default: () => h(icon)
    })
  }
};

const adminOptions = computed(() => {
  let options: DropdownOption[] = [
    {
      label: '删除',
      key: 'delete',
      icon: renderIcon(TrashOutline)
    },
  ];
  if (post.value.is_lock === 0) {
    options.push({
      label: '锁定',
      key: 'lock',
      icon: renderIcon(LockClosedOutline)
    });
  } else {
    options.push({
      label: '解锁',
      key: 'unlock',
      icon: renderIcon(LockOpenOutline)
    });
  }
  if (store.state.userInfo.is_admin) {
    if (post.value.is_top === 0) {
      options.push({
        label: '置顶',
        key: 'stick',
        icon: renderIcon(PushOutline)
      });
    } else {
      options.push({
        label: '取消置顶',
        key: 'unstick',
        icon: renderIcon(PushOutline)
      });
    }
  }
  if (post.value.is_essence === 0) {
    options.push({
      label: '设为亮点',
      key: 'highlight',
      icon: renderIcon(FlameOutline)
    });
  } else {
    options.push({
      label: '取消亮点',
      key: 'unhighlight',
      icon: renderIcon(FlameOutline)
    });
  }
  if (post.value.visibility === VisibilityEnum.PUBLIC) {
    options.push({
      label: '公开',
      key: 'vpublic',
      icon: renderIcon(EyeOutline),
      children: [
        { label: '私密', key: 'vprivate', icon: renderIcon(EyeOffOutline) }
        , { label: '好友可见', key: 'vfriend', icon: renderIcon(PersonOutline) }
      ]
    })
  } else if (post.value.visibility === VisibilityEnum.PRIVATE) {
    options.push({
      label: '私密',
      key: 'vprivate',
      icon: renderIcon(EyeOffOutline),
      children: [
        { label: '公开', key: 'vpublic', icon: renderIcon(EyeOutline) }
        , { label: '好友可见', key: 'vfriend', icon: renderIcon(PersonOutline) }
      ]
    })
  } else {
    options.push({
      label: '好友可见',
      key: 'vfriend',
      icon: renderIcon(PersonOutline),
      children: [
        { label: '公开', key: 'vpublic', icon: renderIcon(EyeOutline) }
        , { label: '私密', key: 'vprivate', icon: renderIcon(EyeOffOutline) }
      ]
    })
  }
  return options;
});

const handlePostAction = (
  item: 'delete' | 'lock' | 'unlock' | 'stick' | 'unstick' | 'highlight' | 'unhighlight' | 'vpublic' | 'vprivate' | 'vfriend'
) => {
  switch (item) {
    case 'delete':
      showDelModal.value = true;
      break;
    case 'lock':
    case 'unlock':
      showLockModal.value = true;
      break;
    case 'stick':
    case 'unstick':
      showStickModal.value = true;
      break;
    case 'highlight':
    case 'unhighlight':
      showHighlightModal.value = true;
      break;
    case 'vpublic':
      tempVisibility.value = 0;
      showVisibilityModal.value = true;
      break;
    case 'vprivate':
      tempVisibility.value = 1;
      showVisibilityModal.value = true;
      break;
    case 'vfriend':
      tempVisibility.value = 2;
      showVisibilityModal.value = true;
      break;
    default:
      break;
  }
};

const post = computed(() => {
  let post: Item.PostComponentProps = Object.assign(
    {
      texts: [],
      imgs: [],
      videos: [],
      links: [],
      attachments: [],
      charge_attachments: [],
    },
    props.post
  );
  post.contents.map((content) => {
    if (+content.type === 1 || +content.type === 2) {
      post.texts.push(content);
    }
    if (+content.type === 3) {
      post.imgs.push(content);
    }
    if (+content.type === 4) {
      post.videos.push(content);
    }
    if (+content.type === 6) {
      post.links.push(content);
    }
    if (+content.type === 7) {
      post.attachments.push(content);
    }
    if (+content.type === 8) {
      post.charge_attachments.push(content);
    }
  });
  return post;
});



const execDelAction = () => {
  deletePost({
    id: post.value.id,
  })
    .then((_res) => {
      window.$message.success('删除成功');
      router.replace('/');

      setTimeout(() => {
        store.commit('refresh');
      }, 50);
    })
    .catch((_err) => {
      loading.value = false;
    });
};
const execLockAction = () => {
  lockPost({
    id: post.value.id,
  })
    .then((res) => {
      emit('reload');
      if (res.lock_status === 1) {
        window.$message.success('锁定成功');
      } else {
        window.$message.success('解锁成功');
      }
    })
    .catch((_err) => {
      loading.value = false;
    });
};
const execStickAction = () => {
  stickPost({
    id: post.value.id,
  })
    .then((res) => {
      emit('reload');
      if (res.top_status === 1) {
        window.$message.success('置顶成功');
      } else {
        window.$message.success('取消置顶成功');
      }
    })
    .catch((_err) => {
      loading.value = false;
    });
};

const execVisibilityAction = () => {
  visibilityPost({
    id: post.value.id,
    visibility: tempVisibility.value
  })
    .then((res) => {
      emit('reload');
      window.$message.success('修改可见性成功');
    })
    .catch((_err) => {
      loading.value = false;
    });
};


const goPostDetail = (id: number) => {
  router.push({
    name: 'post',
    query: {
      id,
    },
  });
};
const doClickText = (e: MouseEvent, id: number) => {
  if ((e.target as any).dataset.detail) {
    const d = (e.target as any).dataset.detail.split(':');
    if (d.length === 2) {
      store.commit('refresh');
      if (d[0] === 'tag') {
        router.push({
          name: 'home',
          query: {
            q: d[1],
            t: 'tag',
          },
        });
      } else {
        router.push({
          name: 'user',
          query: {
            s: d[1],
          },
        });
      }
      return;
    }
  }
  goPostDetail(id);
};
</script>

<style lang="less">
.post-item {
  width: 100%;
  padding: 16px;
  box-sizing: border-box;

  .nickname-wrap {
    font-size: 14px;
  }

  .username-wrap {
    font-size: 14px;
    opacity: 0.75;
  }

  .top-tag {
    transform: scale(0.75);
  }

  .item-header-extra {
    display: flex;
    align-items: center;
    opacity: 0.75;

    .timestamp {
      font-size: 12px;
    }
  }

  .post-text {
    text-align: justify;
    overflow: hidden;
    white-space: pre-wrap;
    word-break: break-all;
  }

  .opt-item {
    display: flex;
    align-items: center;
    opacity: 0.7;

    .opt-item-icon {
      margin-right: 10px;
    }
  }

  &:hover {
    background: #f7f9f9;
    cursor: pointer;
  }

  .n-thing-avatar {
    margin-top: 0;
  }

  .n-thing-header {
    line-height: 16px;
    margin-bottom: 8px !important;
  }
}

.dark {
  .post-item {
    &:hover {
      background: #18181c;
    }

    background-color: rgba(16, 16, 20, 0.75);
  }
}
</style>