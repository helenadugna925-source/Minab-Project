<script setup lang="ts">
const isEventsOpen = ref(false);
const token = useCookie("token");

const navTo = (path: string) => {
  isEventsOpen.value = false; // Close menu first
  navigateTo(path); // Then change page
};
</script>

<template>
  <nav
    class="bg-white border-b border-slate-100 sticky top-0 z-[200] font-normal"
  >
    <div class="max-w-7xl mx-auto flex items-center justify-between p-4">
      <NuxtLink to="/" class="text-lg font-normal tracking-tight"
        >Minab Events</NuxtLink
      >

      <div class="relative">
        <button
          @click="isEventsOpen = !isEventsOpen"
          class="px-4 py-2 text-sm text-slate-600"
        >
          Events ▼
        </button>
        <div
          v-if="isEventsOpen"
          class="absolute right-0 mt-2 w-52 bg-white rounded-2xl shadow-2xl border border-slate-50 py-2 z-[250]"
        >
          <button
            @click="navTo('/events')"
            class="w-full text-left px-6 py-3 text-sm text-slate-500 hover:bg-blue-50 hover:text-blue-600"
          >
            Browse Events
          </button>
          <button
            @click="navTo('/events/my')"
            class="w-full text-left px-6 py-3 text-sm text-slate-500 hover:bg-indigo-50 hover:text-indigo-600 border-t border-slate-50"
          >
            My Events
          </button>
          <button
            @click="
              () => {
                isEventsOpen = false;
                navigateTo('/events/bookmarks');
              }
            "
            class="w-full text-left px-5 py-3 text-sm font-normal text-slate-500 hover:bg-amber-50 hover:text-amber-600 transition-colors border-t border-slate-50"
          >
            Bookmarks
          </button>
          <button
            @click="navTo('/events/reserved')"
            class="w-full text-left px-6 py-3 text-sm text-slate-500 hover:bg-emerald-50 hover:text-emerald-600 border-t border-slate-50"
          >
            Reserved Events
          </button>
        </div>
      </div>
    </div>
  </nav>
</template>
