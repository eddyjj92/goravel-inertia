<template>
  <img :src="logo" :width="size" :height="size" alt="Goravel Inertia" class="gi-logo" />
</template>

<script setup lang="ts">
import logo from '../assets/goravel-inertia.png'

withDefaults(defineProps<{ size?: number | string }>(), { size: 40 })
</script>

<style scoped>
.gi-logo {
  display: block;
  object-fit: contain;
}
</style>
