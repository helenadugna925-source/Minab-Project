<script setup lang="ts">
// Define what data the card needs
defineProps<{
  event: {
    id: number | string;
    title: string;
    featured_image: string;
    event_date: string;
    venue_name: string;
    category?: { name: string };
  };
}>();

const fallback =
  "https://images.unsplash.com/photo-1501281668745-f7f57925c3b4?auto=format&fit=crop&q=80&w=800";
</script>

<template>
  <div
    class="group bg-white rounded-[2.5rem] border border-slate-100 overflow-hidden shadow-sm hover:shadow-xl transition-all duration-500 flex flex-col h-full"
  >
    <!-- Image Area -->
    <div class="relative h-56 overflow-hidden bg-slate-100">
      <img
        :src="event.featured_image || fallback"
        class="w-full h-full object-cover group-hover:scale-105 transition-transform duration-1000"
      />
      <div class="absolute top-4 left-4">
        <span
          class="px-3 py-1 bg-white/90 backdrop-blur-md rounded-xl text-[9px] font-normal text-indigo-600 uppercase tracking-widest border border-white/50 shadow-sm"
        >
          {{ event.category?.name || "General" }}
        </span>
      </div>
    </div>

    <!-- Content Area -->
    <div class="p-8 flex flex-col flex-grow">
      <p
        class="text-indigo-600 text-[10px] uppercase tracking-widest mb-2 font-normal"
      >
        {{ new Date(event.event_date).toDateString() }}
      </p>

      <h3
        class="text-lg font-normal text-slate-900 mb-6 capitalize leading-tight flex-grow"
      >
        {{ event.title }}
      </h3>

      <div class="flex items-center gap-2 text-slate-400 text-xs mb-8 italic">
        <span>📍</span> {{ event.venue_name }}
      </div>

      <!-- THE MAGIC BUTTON (Works everywhere) -->
      <NuxtLink
        :to="`/events/${event.id}`"
        class="block w-full py-4 bg-slate-900 text-white text-center rounded-2xl text-[10px] font-normal uppercase tracking-widest hover:bg-indigo-600 transition-all flex items-center justify-center gap-2 group/btn cursor-pointer"
      >
        View Details
        <span class="group-hover/btn:translate-x-1 transition-transform"
          >→</span
        >
      </NuxtLink>
    </div>
  </div>
</template>

<style scoped>
/* Ensuring normal weight as requested */
* {
  font-weight: 400 !important;
}
</style>
