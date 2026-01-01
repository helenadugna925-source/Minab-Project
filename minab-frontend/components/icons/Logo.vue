<template>
  <!-- Logo graphic: outer color controlled by CSS via `currentColor` -->
  <svg
    v-bind="$attrs"
    xmlns="http://www.w3.org/2000/svg"
    viewBox="0 0 24 24"
    role="img"
    aria-label="Minab Events logo"
  >
    <!-- background shape uses currentColor so class `text-primary` or `text-accent` will color it -->
    <rect width="24" height="24" rx="5" fill="currentColor" />
    <!-- inner mark in white for contrast -->
    <path fill="#ffffff" d="M6.5 16.5l2.5-6 2 3 2-4 2.5 7.5H6.5z" />
  </svg>
</template>

<script setup lang="ts">
// Simple stateless logo component that passes attributes (class/size) to the root svg
</script>
