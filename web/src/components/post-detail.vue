<template>
  <div class="detail-item" @click="goPostDetail(post.id)">
    <n-thing>
      <template #avatar>
        <n-avatar round :size="30" :src="post.user.avatar" />
      </template>
      <template #header>
        <router-link @click.stop class="username-link" :to="{
          name: 'user',
          query: { s: post.user.username },
        }">
          {{ post.user.nickname }}
        </router-link>
        <span class="username-wrap"> @{{ post.user.username }} </span>
        <n-tag v-if="post.is_top" class="top-tag" type="warning" size="small" round>
          置顶
        </n-tag>
        <n-tag v-if="post.visibility == VisibilityEnum.PRIVATE" class="top-tag" type="error" size="small" round>
          私密
        </n-tag>
        <n-tag v-if="post.visibility == VisibilityEnum.FRIEND" class="top-tag" type="info" size="small" round>
          好友可见
        </n-tag>
      </template>
      <template #header-extra>
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


        <!-- 编辑content -->
        <n-modal v-model:show="showEditModal" :mask-closable="false" class="custom-card" style="width:600px;"
          preset="card" title="修改" positive-text="确认" negative-text="取消" @positive-click="execEditAction">
          <template #header-extra>
          </template>
          <n-input type="textarea" v-model:value="postContentText.content" />
          <template #footer>
            <n-button @click="execEditAction" type="primary" style="float:right;">
              保存
            </n-button>
          </template>
        </n-modal>

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
        <!-- 亮点确认 -->
        <n-modal v-model:show="showHighlightModal" :mask-closable="false" preset="dialog" title="提示" :content="'确定将该泡泡动态' +
          (post.is_essence ? '取消亮点' : '设为亮点') +
          '吗？'
          " positive-text="确认" negative-text="取消" @positive-click="execHighlightAction" />
        <!-- 修改可见度确认 -->
        <n-modal v-model:show="showVisibilityModal" :mask-closable="false" preset="dialog" title="提示" :content="'确定将该泡泡动态可见度修改为' +
          (tempVisibility == 0 ? '公开' : (tempVisibility == 1 ? '私密' : '好友可见')) +
          '吗？'
          " positive-text="确认" negative-text="取消" @positive-click="execVisibilityAction" />
      </template>
      <div v-if="post.texts.length > 0">
        <span v-for="content in post.texts" :key="content.id" class="post-text" @click.stop="doClickText($event, post.id)"
          v-html="parsePostTag(content.content).content">
        </span>
      </div>

      <template #footer>
        <post-attachment :attachments="post.attachments" />
        <post-attachment :attachments="post.charge_attachments" :price="post.attachment_price" />
        <post-image :imgs="post.imgs" />
        <post-video :videos="post.videos" :full="true" />
        <post-link :links="post.links" />
        <div class="timestamp">
          发布于 {{ formatPrettyTime(post.created_on) }}
          <span v-if="post.ip_loc">
            <n-divider vertical />
            {{ post.ip_loc }}
          </span>
          <span v-if="!store.state.collapsedLeft && post.created_on != post.latest_replied_on">
            <n-divider vertical /> 最后回复
            {{ formatPrettyTime(post.latest_replied_on) }}
          </span>
        </div>
      </template>
      <template #action>
        <div class="opts-wrap">
          <n-space justify="space-between">
            <div class="opt-item hover" @click.stop="handlePostStar">
              <n-icon size="20" class="opt-item-icon">
                <heart-outline v-if="!hasStarred" />
                <heart v-if="hasStarred" color="red" />
              </n-icon>
              {{ post.upvote_count }}
            </div>
            <div class="opt-item">
              <n-icon size="20" class="opt-item-icon">
                <chatbox-outline />
              </n-icon>
              {{ post.comment_count }}
            </div>
            <div class="opt-item hover" @click.stop="handlePostCollection">
              <n-icon size="20" class="opt-item-icon">
                <bookmark-outline v-if="!hasCollected" />
                <bookmark v-if="hasCollected" color="#ff7600" />
              </n-icon>
              {{ post.collection_count }}
            </div>
            <div class="opt-item hover" @click.stop="handlePostShare">
              <n-icon size="20" class="opt-item-icon">
                <share-social-outline />
              </n-icon>
              {{ post.share_count }}
            </div>
          </n-space>
        </div>
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
import { formatPrettyTime } from '@/utils/formatTime';
import { parsePostTag } from '@/utils/content';
import {
  Heart,
  HeartOutline,
  Bookmark,
  BookmarkOutline,
  ShareSocialOutline,
  ChatboxOutline,
  PushOutline,
  TrashOutline,
  LockClosedOutline,
  CreateOutline,
  LockOpenOutline,
  EyeOutline,
  EyeOffOutline,
  PersonOutline,
  FlameOutline,
} from '@vicons/ionicons5';
import { MoreHorizFilled } from '@vicons/material';
import {
  getPostStar,
  postStar,
  getPostCollection,
  postCollection,
  editPostText,
  deletePost,
  lockPost,
  stickPost,
  highlightPost,
  visibilityPost
} from '@/api/post';
import type { DropdownOption } from 'naive-ui';
import { VisibilityEnum } from '@/utils/IEnum';
import copy from "copy-to-clipboard";

const store = useStore();
const router = useRouter();
const hasStarred = ref(false);
const hasCollected = ref(false);
const props = withDefaults(
  defineProps<{
    post: Item.PostProps;
  }>(),
  {}
);
const showEditModal = ref(false);
const showDelModal = ref(false);
const showLockModal = ref(false);
const showStickModal = ref(false);
const showHighlightModal = ref(false);
const showVisibilityModal = ref(false);
const loading = ref(false);
const tempVisibility = ref<VisibilityEnum>(VisibilityEnum.PUBLIC);

// cyw 内容编辑
const postContentText = ref({
  id: 0,
  post_id: 0,
  user_id: 0,
  type: 0,
  sort: 0,
  content: ""
})

const emit = defineEmits<{
  (e: 'reload'): void;
}>();

const post = computed({
  get: () => {
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


    for (let i = 0; i < post.contents.length; i++) {
      const element = post.contents[i];
      if (element.type === 2) {
        postContentText.value = {
          id: element.id,
          content: element.content,
          post_id: element.post_id,
          user_id: post.user_id || 0,
          type: element.type,
          sort: element.sort
        }
      }
    }

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
  },
  set: (newVal) => {
    props.post.upvote_count = newVal.upvote_count;
    props.post.comment_count = newVal.comment_count;
    props.post.collection_count = newVal.collection_count;
    props.post.is_essence = newVal.is_essence;
  },
});

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

  // 编辑
  options.push({
    label: '编辑',
    key: 'edit',
    icon: renderIcon(CreateOutline)
  });

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
const handlePostAction = (
  item: 'delete' | 'lock' | 'unlock' | 'stick' | 'unstick' | 'highlight' | 'unhighlight' | 'vpublic' | 'vprivate' | 'vfriend' | 'edit'
) => {
  switch (item) {
    case 'edit':
      showEditModal.value = true;
      break;
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

// 编辑
const execEditAction = () => {
  editPostText(postContentText.value)
    .then((_res) => {
      window.$message.success('修改成功');
      showEditModal.value = false;
      setTimeout(() => {
        store.commit('refresh');
      }, 50);
    })
    .catch((_err) => {
      loading.value = false;
    });
};

// 删除
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
const execHighlightAction = () => {
  highlightPost({
    id: post.value.id,
  })
    .then((res) => {
      post.value = {
        ...post.value,
        is_essence: res.highlight_status,
      };
      if (res.highlight_status === 1) {
        window.$message.success('设为亮点成功');
      } else {
        window.$message.success('取消亮点成功');
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
const handlePostStar = () => {
  postStar({
    id: post.value.id,
  })
    .then((res) => {
      hasStarred.value = res.status;
      if (res.status) {
        post.value = {
          ...post.value,
          upvote_count: post.value.upvote_count + 1,
        };
      } else {
        post.value = {
          ...post.value,
          upvote_count: post.value.upvote_count - 1,
        };
      }
    })
    .catch((err) => {
      console.log(err);
    });
};
const handlePostCollection = () => {
  postCollection({
    id: post.value.id,
  })
    .then((res) => {
      hasCollected.value = res.status;
      if (res.status) {
        post.value = {
          ...post.value,
          collection_count: post.value.collection_count + 1,
        };
      } else {
        post.value = {
          ...post.value,
          collection_count: post.value.collection_count - 1,
        };
      }
    })
    .catch((err) => {
      console.log(err);
    });
};
const handlePostShare = () => {
  copy(`${window.location.origin}/#/post?id=${post.value.id}`);
  window.$message.success('链接已复制到剪贴板');
};

onMounted(() => {
  if (store.state.userInfo.id > 0) {
    getPostStar({
      id: post.value.id,
    })
      .then((res) => {
        hasStarred.value = res.status;
      })
      .catch((err) => {
        console.log(err);
      });

    getPostCollection({
      id: post.value.id,
    })
      .then((res) => {
        hasCollected.value = res.status;
      })
      .catch((err) => {
        console.log(err);
      });

    console.log(post.value)
  }
});
</script>

<style lang="less">
.detail-item {
  width: 100%;
  padding: 16px;
  box-sizing: border-box;

  background: #f7f9f9;

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

  .options {
    opacity: 0.75;
  }

  .post-text {
    font-size: 16px;
    text-align: justify;
    overflow: hidden;
    white-space: pre-wrap;
    word-break: break-all;
  }

  .opts-wrap {
    margin-top: 20px;

    .opt-item {
      display: flex;
      align-items: center;
      opacity: 0.7;

      .opt-item-icon {
        margin-right: 10px;
      }

      &.hover {
        cursor: pointer;
      }
    }
  }

  .n-thing {
    .n-thing-avatar-header-wrapper {
      align-items: center;
    }
  }

  .timestamp {
    opacity: 0.75;
    font-size: 12px;
    margin-top: 10px;
  }
}

.dark {
  .detail-item {
    background: #18181c;
  }
}
</style>