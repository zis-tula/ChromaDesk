<template>
  <div class="list-editor">
    <div v-for="(item, index) in modelValue" :key="index" class="list-editor-row">
      <el-input
        :model-value="item"
        @update:model-value="updateItem(index, $event)"
        :placeholder="placeholder"
        size="small"
        class="list-editor-input"
      />
      <el-button
        type="danger"
        size="small"
        text
        :icon="Delete"
        @click="removeItem(index)"
        class="list-editor-delete"
      />
    </div>
    <el-button type="primary" size="small" text :icon="Plus" @click="addItem">
      {{ addText }}
    </el-button>
  </div>
</template>

<script setup>
import { Plus, Delete } from '@element-plus/icons-vue'

const props = defineProps({
  modelValue: { type: Array, default: () => [] },
  placeholder: { type: String, default: '' },
  addText: { type: String, default: 'Add' }
})

const emit = defineEmits(['update:modelValue'])

function addItem() {
  emit('update:modelValue', [...props.modelValue, ''])
}

function removeItem(index) {
  const list = [...props.modelValue]
  list.splice(index, 1)
  emit('update:modelValue', list)
}

function updateItem(index, value) {
  const list = [...props.modelValue]
  list[index] = value
  emit('update:modelValue', list)
}
</script>

<style scoped>
.list-editor {
  width: 100%;
  max-width: 500px;
}

.list-editor-row {
  display: flex;
  align-items: center;
  gap: 4px;
  margin-bottom: 8px;
}

.list-editor-input {
  flex: 1;
}

.list-editor-delete {
  flex-shrink: 0;
}
</style>
