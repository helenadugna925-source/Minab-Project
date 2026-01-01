<script setup lang="ts">
import { ref, onMounted } from "vue";
import { jwtDecode } from "jwt-decode";
import { INSERT_TICKET_MUTATION } from "~/utils/constants/queries/events";

const isFinished = ref(false);
const isLoading = ref(true);
const token = useCookie("token");

// 1. Setup the mutation
const { mutate: saveTicket } = useMutation(INSERT_TICKET_MUTATION);

onMounted(async () => {
  // 2. Get Event ID from memory
  const eventId = localStorage.getItem("pending_event_id");

  // 3. Get your REAL user ID from token
  let realUserId = 1;
  if (token.value) {
    try {
      const decoded: any = jwtDecode(token.value);
      realUserId =
        decoded["https://hasura.io/jwt/claims"]?.["x-hasura-user-id"] ||
        decoded.sub ||
        1;
    } catch (e) {}
  }

  try {
    // 4. THE FIX: Manually insert ticket to Database
    if (eventId) {
      await saveTicket({
        event_id: parseInt(eventId),
        user_id: parseInt(realUserId.toString()),
        ticket_number: "TICK-" + Math.floor(Math.random() * 1000000),
      });
      localStorage.removeItem("pending_event_id"); // Clear memory
    }
    setTimeout(() => {
      isFinished.value = true;
      isLoading.value = false;
    }, 2000);
  } catch (err) {
    console.error("Save Error:", err);
    isLoading.value = false;
    isFinished.value = true;
  }
});
</script>

<template>
  <div
    class="min-h-screen bg-white flex flex-col items-center justify-center p-6 text-center font-normal"
  >
    <div v-if="isLoading">
      <div
        class="w-12 h-12 border-2 border-indigo-600 border-t-transparent rounded-full animate-spin mx-auto mb-4"
      ></div>
      <p class="text-slate-400 text-[10px] uppercase tracking-[0.3em]">
        Confirming your ticket...
      </p>
    </div>
    <div v-else-if="isFinished" class="max-w-md w-full">
      <div
        class="bg-white rounded-[3rem] border border-slate-50 shadow-2xl p-12"
      >
        <div
          class="w-16 h-16 bg-emerald-50 text-emerald-500 rounded-full flex items-center justify-center mx-auto mb-8 shadow-inner text-xl"
        >
          ✓
        </div>
        <h1
          class="text-3xl font-normal text-slate-900 mb-4 uppercase tracking-tighter"
        >
          Success!
        </h1>
        <p class="text-slate-500 font-normal text-sm mb-10 italic">
          "Your digital ticket is ready. You can find it in the Reserved Events
          section."
        </p>
        <button
          @click="navigateTo('/events/reserved')"
          class="w-full py-5 bg-indigo-600 text-white rounded-2xl text-[10px] font-normal uppercase tracking-widest shadow-lg shadow-indigo-100 transition-all"
        >
          View My Ticket →
        </button>
      </div>
    </div>
  </div>
</template>
