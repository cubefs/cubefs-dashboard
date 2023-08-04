<!--
 Copyright 2023 The CubeFS Authors.

 Licensed under the Apache License, Version 2.0 (the "License");
 you may not use this file except in compliance with the License.
 You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or
 implied. See the License for the specific language governing
 permissions and limitations under the License.
-->

<script>
const renderLi = (h, list) => {
  return list.map(item => {
    return h('li', {}, [item])
  })
}
export default {
  functional: true,
  props: {
    count: {
      type: Number,
      default: 4
    },
    title: {
      type: String,
      default: '更多操作'
    },
    type: {
      type: String,
      default: 'text'
    },
    poperClass: {
      type: String,
      default: ''
    },
    width: {
      type: [String, Number],
      default: 'auto'
    },
    visibleArrow: {
      type: Boolean,
      default: true
    }
  },
  render(h, context) {
    const { count, title, type, poperClass, width, visibleArrow } = context.props
    context.children = context.children.filter(item => {
      return item.tag
    })
    const children = context.children
    const childrenLen = children.length
    return h(
      'div',
      {},
      [
        ...children.slice(0, childrenLen > count ? count - 1 : count),
        childrenLen > count
          ? h('o-dropdown', {
            class: 'm-l o-dropdown-arrow-text',
            props: {
              type, poperClass, width, visibleArrow
            },
            scopedSlots: {
              header: props => h('span', title),
              list: props =>
                h('ul', {}, [...renderLi(h, children.slice(count - 1))])
            }
          })
          : ''
      ]
    )
  }
}
</script>
<style lang="scss" scoped>
.m-l {
  margin-left: 5px;
}
</style>
