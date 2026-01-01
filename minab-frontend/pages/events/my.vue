<script setup lang="ts">
import { ref, onMounted, computed } from "vue";
import { jwtDecode } from "jwt-decode";
import { GetMyEventsQuery } from "~/utils/constants/queries/events";

const token = useCookie("token");
const myEvents = ref<any[]>([]);
const totalCount = ref(0);
const isLoading = ref(true);

// GET REAL ID FROM TOKEN
const myUserId = computed(() => {
  if (!token.value) return 1;
  try {
    const decoded: any = jwtDecode(token.value);
    return parseInt(
      decoded["https://hasura.io/jwt/claims"]?.["x-hasura-user-id"] ||
        decoded.sub ||
        "1"
    );
  } catch (e) {
    return 1;
  }
});

// FETCH DATA using the real ID
const { data, refresh } = await useAsyncQuery<any>(GetMyEventsQuery, {
  user_id: myUserId.value,
  skip: 0,
  take: 10,
});

onMounted(async () => {
  isLoading.value = true;
  await refresh(); // Force fresh data from Postgres
  if (data.value) {
    myEvents.value = data.value.events || [];
    totalCount.value = data.value.events_aggregate?.aggregate?.count || 0;
  }
  isLoading.value = false;
});

const fallback =
  "https://images.unsplash.com/photo-1501281668745-f7f57925c3b4?auto=format&fit=crop&q=80&w=800";
</script>

<template>
  <div class="min-h-screen bg-slate-50/50 pb-24 font-normal relative">
    <header
      class="bg-white border-b py-10 px-8 flex justify-between items-center"
    >
      <h1
        class="text-3xl font-normal uppercase tracking-tighter text-slate-900"
      >
        My Events
      </h1>
      <div
        class="bg-indigo-600 text-white px-6 py-2 rounded-2xl font-normal shadow-lg text-xs uppercase"
      >
        Total Hosted: {{ totalCount }}
      </div>
    </header>

    <main
      class="max-w-7xl mx-auto px-8 mt-12 grid grid-cols-1 md:grid-cols-3 gap-10"
    >
      <div
        v-if="myEvents.length === 0"
        class="col-span-full text-center py-20 text-slate-400"
      >
        No events found for your account.
      </div>

      <div
        v-for="event in myEvents"
        :key="event.id"
        class="group bg-white rounded-[2.5rem] border border-slate-100 overflow-hidden shadow-sm hover:shadow-xl transition-all duration-500 flex flex-col h-full"
      >
        <img
          :src="event.featured_image || fallback"
          class="h-60 w-full object-cover"
        />
        <div class="p-8 flex flex-col flex-grow">
          <h3
            class="text-xl font-normal text-slate-900 mb-6 capitalize leading-tight"
          >
            {{ event.title }}
          </h3>
          <NuxtLink
            :to="'/events/' + event.id"
            class="mt-auto block w-full py-4 bg-slate-900 text-white text-center rounded-2xl text-[10px] uppercase tracking-widest hover:bg-indigo-600 transition-all font-normal"
            >View Details →</NuxtLink
          >
        </div>
      </div>
    </main>

    <NuxtLink
      to="/events/create"
      class="fixed bottom-10 right-10 w-16 h-16 bg-indigo-600 text-white rounded-full flex items-center justify-center shadow-2xl hover:scale-110 transition-all z-50 group"
    >
      <svg
        xmlns="http://www.w3.org/2000/svg"
        class="w-8 h-8 group-hover:rotate-90 transition-transform duration-300"
        fill="none"
        viewBox="0 0 24 24"
        stroke="currentColor"
        stroke-width="2.5"
      >
        <path d="M12 4v16m8-8H4" />
      </svg>
    </NuxtLink>
  </div>
</template>

<style scoped>
* {
  font-weight: 400 !important;
}
</style>
