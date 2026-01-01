<script setup lang="ts">
import { ref, computed, onMounted } from "vue";
import { jwtDecode } from "jwt-decode";
import {
  GetEventByIdQuery,
  INITIALIZE_PAYMENT_MUTATION,
  INSERT_BOOKMARK_MUTATION,
  DELETE_BOOKMARK_MUTATION,
} from "~/utils/constants/queries/events";

const route = useRoute();
const eventId = computed(() => parseInt(route.params.id as string));
const isBooking = ref(false);
const isBookmarked = ref(false);

const token = useCookie("token");
const decoded = ref<any>({});
const myId = ref(1);

if (token.value) {
  try {
    decoded.value = jwtDecode(token.value);
    const idVal =
      decoded.value["https://hasura.io/jwt/claims"]?.["x-hasura-user-id"] ||
      decoded.value.sub;
    myId.value = parseInt(idVal || "1");
  } catch (e) {}
}

const { data, refresh } = await useAsyncQuery<any>(GetEventByIdQuery, {
  id: eventId.value,
  user_id: myId.value,
});
const event = computed(() => data.value?.events_by_pk);

onMounted(() => {
  if (event.value?.event_bookmarks?.length > 0) isBookmarked.value = true;
});

const { mutate: pay } = useMutation(INITIALIZE_PAYMENT_MUTATION);
const handleBooking = async () => {
  localStorage.setItem("pending_event_id", eventId.value.toString());
  isBooking.value = true;
  try {
    const res = await pay({
      event_id: eventId.value,
      full_name: decoded.value?.name || "Customer",
      email: decoded.value?.email || "customer@gmail.com",
    });
    if (res?.data?.initializePayment?.checkout_url)
      window.location.href = res.data.initializePayment.checkout_url;
  } catch (err: any) {
    alert("Payment Error: " + err.message);
  } finally {
    isBooking.value = false;
  }
};

// 2. BOOKMARK TOGGLE LOGIC (MARK/UNMARK)
const { mutate: saveB } = useMutation(INSERT_BOOKMARK_MUTATION);
const { mutate: delB } = useMutation(DELETE_BOOKMARK_MUTATION);

const toggleBookmark = async () => {
  try {
    if (isBookmarked.value) {
      await delB({ event_id: eventId.value, user_id: myId.value });
      isBookmarked.value = false;
    } else {
      await saveB({ event_id: eventId.value, user_id: myId.value });
      isBookmarked.value = true;
    }
    await refresh();
  } catch (e) {
    alert("Please log in again.");
  }
};
</script>

<template>
  <div v-if="event" class="min-h-screen bg-white font-normal pb-20">
    <div class="relative h-[45vh] bg-slate-900 flex items-end p-12 text-white">
      <h1 class="text-6xl tracking-tighter uppercase">{{ event.title }}</h1>
    </div>
    <main class="max-w-7xl mx-auto grid lg:grid-cols-12 gap-16 py-16 px-6">
      <div class="lg:col-span-8">
        <p class="text-slate-800 text-4xl border-l-8 border-indigo-100 pl-10">
          "{{ event.description }}"
        </p>
      </div>
      <div class="lg:col-span-4">
        <div
          class="sticky top-28 bg-white border p-10 rounded-[3rem] shadow-2xl text-center space-y-8"
        >
          <span class="text-5xl text-slate-900 tracking-tighter"
            >{{ event.price }} ETB</span
          >
          <div class="space-y-4">
            <button
              @click="handleBooking"
              :disabled="isBooking"
              class="w-full py-6 bg-indigo-600 text-white rounded-[2rem] text-sm uppercase shadow-xl active:scale-95 transition-all"
            >
              Reserve Spot Now →
            </button>
            <button
              @click="toggleBookmark"
              class="w-full py-5 border border-slate-100 rounded-[2rem] text-[10px] uppercase transition-all"
              :class="
                isBookmarked
                  ? 'bg-amber-50 text-amber-600 border-amber-200 shadow-inner'
                  : 'text-slate-400'
              "
            >
              {{ isBookmarked ? "★ Saved" : "☆ Bookmark Event" }}
            </button>
          </div>
        </div>
      </div>
    </main>
  </div>
</template>
<style scoped>
* {
  font-weight: 400 !important;
}
</style>
