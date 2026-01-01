<script setup lang="ts">
import { ref, computed, onMounted } from "vue";
import { jwtDecode } from "jwt-decode";
import { GetReservedEventsQuery } from "~/utils/constants/queries/events";

const token = useCookie("token");
const userId = ref(1);

if (token.value) {
  try {
    const decoded: any = jwtDecode(token.value);
    userId.value =
      decoded["https://hasura.io/jwt/claims"]?.["x-hasura-user-id"] ||
      decoded.sub ||
      1;
  } catch (e) {}
}

const { data, refresh } = await useAsyncQuery<any>(GetReservedEventsQuery, {
  user_id: parseInt(userId.value.toString()),
});

onMounted(() => refresh());
</script>

<template>
  <div class="min-h-screen bg-slate-50/50 p-10 font-normal">
    <div class="max-w-5xl mx-auto">
      <h1 class="text-2xl uppercase tracking-widest mb-10">
        My Confirmed Tickets
      </h1>
      <div
        v-if="!data?.tickets?.length"
        class="text-center py-20 text-slate-400"
      >
        No tickets found.
      </div>
      <div class="grid grid-cols-1 md:grid-cols-2 gap-8">
        <div
          v-for="res in data?.tickets"
          :key="res.id"
          class="flex bg-white rounded-[2rem] overflow-hidden border border-slate-100 shadow-sm relative h-44 group"
        >
          <div class="w-1/3 bg-slate-900">
            <img
              :src="res.event?.featured_image"
              class="w-full h-full object-cover opacity-60 group-hover:scale-105 transition-all duration-1000"
            />
          </div>
          <div class="w-2/3 p-8 flex flex-col justify-between">
            <h3 class="text-lg leading-tight mb-1 capitalize">
              {{ res.event?.title }}
            </h3>
            <p class="text-[9px] text-slate-400 uppercase">
              ID: {{ res.ticket_number }}
            </p>
            <!-- FIXED VIEW DETAILS LINK -->
            <NuxtLink
              :to="'/events/' + res.event?.id"
              class="text-indigo-600 text-[10px] font-normal uppercase underline tracking-widest"
              >Details →</NuxtLink
            >
          </div>
          <div
            class="absolute left-1/3 top-0 bottom-0 w-px border-l-2 border-dashed border-slate-100"
          ></div>
        </div>
      </div>
    </div>
  </div>
</template>
