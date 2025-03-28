<template>
  <div class="exercise-list">
    <a-list :data-source="exercises" size="small">
      <template #renderItem="{ item }">
        <a-list-item>
          <a-list-item-meta>
            <template #title>
              <a :href="item.url" target="_blank">
                {{ item.title }}
              </a>
            </template>
            <template #description>
              <a-tag :color="getDifficultyColor(item.difficulty)">
                {{ item.difficulty }}
              </a-tag>
            </template>
          </a-list-item-meta>
        </a-list-item>
      </template>
    </a-list>
  </div>
</template>

<script setup>
const props = defineProps({
  exercises: {
    type: Array,
    required: true,
    default: () => []
  }
})

const getDifficultyColor = (difficulty) => {
  const colors = {
    easy: 'success',
    medium: 'warning',
    hard: 'error'
  }
  return colors[difficulty] || 'default'
}
</script>

<style scoped>
.exercise-list {
  margin-bottom: 16px;
}
</style> 