<script setup lang="ts">
import { ref, onMounted } from "vue";
import { searchEventQuery } from "~/utils/constants/queries/events";
import FilterIcon from "~/components/icons/Filter.vue";

const text = ref("");
const maxPrice = ref(100000);
const radiusKm = ref(0);
const showFilters = ref(false);
const events = ref<any[]>([]);

const centerLat = 9.0122;
const centerLng = 38.7578;

const performSearch = async () => {
  // 1. Calculate Bounding Box for Radius
  let spatial = { min_lat: -90, max_lat: 90, min_lng: -180, max_lng: 180 };
  if (radiusKm.value > 0) {
    const d = radiusKm.value / 111.1; // 1 degree ~ 111km
    spatial = {
      min_lat: centerLat - d,
      max_lat: centerLat + d,
      min_lng: centerLng - d,
      max_lng: centerLng + d,
    };
  }

  const { data, refresh } = await useAsyncQuery<any>(searchEventQuery, {
    text: text.value ? `%${text.value}%` : "%%",
    max_price: parseFloat(maxPrice.value.toString()),
    ...spatial,
  });
  await refresh();
  events.value = data.value?.events || [];
};

onMounted(() => performSearch());

const fallback =
  "https://images.unsplash.com/photo-1514525253361-bee8718a74a2?q=80&w=1000";
const onImgError = (e: any) => {
  e.target.src = fallback;
};
</script>

<template>
  <div class="min-h-screen bg-slate-50/50 pb-24 font-normal">
    <!-- SEARCH & RADIUS ICON -->
    <div class="max-w-xl mx-auto pt-16 px-6 flex gap-4">
      <div class="relative flex-1">
        <input
          v-model="text"
          @input="performSearch"
          placeholder="Type to search events..."
          class="w-full p-5 rounded-[2rem] border border-slate-200 shadow-sm outline-none focus:ring-2 focus:ring-indigo-500 font-normal text-sm"
        />
      </div>
      <button
        @click="showFilters = true"
        class="p-5 bg-white border border-slate-200 rounded-2xl shadow-sm hover:bg-slate-50 transition-all"
      >
        <FilterIcon class="w-6 h-6 text-slate-400" />
      </button>
    </div>

    <main
      class="max-w-7xl mx-auto px-8 mt-16 grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-10"
    >
      <div
        v-for="event in events"
        :key="event.id"
        class="group bg-white rounded-[2.5rem] border border-slate-100 overflow-hidden shadow-sm hover:shadow-xl transition-all duration-500 flex flex-col h-full"
      >
        <div class="h-60 overflow-hidden relative bg-slate-100">
          <img
            :src="event.featured_image || fallback"
            @error="onImgError"
            class="w-full h-full object-cover group-hover:scale-105 transition-transform duration-1000"
          />
          <div
            class="absolute top-4 left-4 bg-white/90 px-3 py-1 rounded-lg text-[9px] text-indigo-600 uppercase tracking-widest font-normal"
          >
            {{ event.category?.name || "General" }}
          </div>
        </div>
        <div class="p-8 flex flex-col flex-grow">
          <h3
            class="text-lg font-normal text-slate-900 mb-6 leading-tight capitalize"
          >
            {{ event.title }}
          </h3>
          <!-- FIXED LINK -->
          <NuxtLink
            :to="'/events/' + event.id"
            class="mt-auto block w-full py-4 bg-slate-900 text-white text-center rounded-2xl text-[10px] uppercase tracking-widest hover:bg-indigo-600 transition-all font-normal"
          >
            View Details →
          </NuxtLink>
        </div>
      </div>
    </main>

    <div
      v-if="showFilters"
      class="fixed inset-0 z-[600] bg-slate-900/40 backdrop-blur-sm flex items-center justify-center p-6"
    >
      <div
        class="bg-white rounded-[3.5rem] shadow-2xl w-full max-w-md p-10 relative text-slate-900"
      >
        <button
          @click="showFilters = false"
          class="absolute top-10 right-10 text-slate-300 hover:text-slate-900 transition-colors"
        >
          ✕
        </button>
        <h2 class="text-2xl mb-10 uppercase tracking-tighter text-center">
          Refine Search
        </h2>
        <div class="space-y-10">
          <div>
            <label
              class="text-[10px] uppercase tracking-widest text-slate-400 mb-4 block"
              >Max Price: {{ maxPrice }} ETB</label
            >
            <input
              type="range"
              min="0"
              max="100000"
              step="500"
              v-model="maxPrice"
              @change="performSearch"
              class="w-full accent-indigo-600 cursor-pointer"
            />
          </div>
          <div>
            <label
              class="text-[10px] uppercase tracking-widest text-slate-400 mb-4 block"
              >Search Radius:
              {{ radiusKm == 0 ? "Unlimited" : radiusKm + " KM" }}</label
            >
            <input
              type="range"
              min="0"
              max="500"
              v-model="radiusKm"
              @change="performSearch"
              class="w-full accent-indigo-600 cursor-pointer"
            />
          </div>
          <button
            @click="showFilters = false"
            class="w-full py-5 bg-indigo-600 text-white rounded-[1.5rem] font-normal text-xs uppercase tracking-widest active:scale-95 transition-all"
          >
            Apply Filters
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
* {
  font-weight: 400 !important;
}
</style>
