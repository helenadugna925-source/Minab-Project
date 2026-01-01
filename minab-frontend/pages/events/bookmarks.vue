<script setup lang="ts">
import {
  GetBookmarkedEventsQuery,
  DELETE_BOOKMARK_MUTATION,
} from "~/utils/constants/queries/events";
import { jwtDecode } from "jwt-decode";

const token = useCookie("token");
const myId = ref(1);
if (token.value) {
  try {
    myId.value = parseInt(
      jwtDecode<any>(token.value!)["https://hasura.io/jwt/claims"][
        "x-hasura-user-id"
      ]
    );
  } catch (e) {}
}

const { data, refresh } = await useAsyncQuery<any>(GetBookmarkedEventsQuery, {
  user_id: myId.value,
});
const { mutate: delB } = useMutation(DELETE_BOOKMARK_MUTATION);

const handleRemove = async (eventId: number) => {
  try {
    await delB({ event_id: eventId, user_id: myId.value });
    await refresh(); // Makes the card disappear instantly!
  } catch (e) {
    alert("Could not remove.");
  }
};
</script>

<template>
  <div class="min-h-screen bg-slate-50/50 p-10 font-normal">
    <h1 class="text-3xl uppercase mb-12 tracking-widest">My Bookmarks</h1>
    <div class="grid grid-cols-3 gap-10">
      <div
        v-for="item in data?.event_bookmarks"
        :key="item.event.id"
        class="bg-white rounded-[2rem] border overflow-hidden flex flex-col shadow-sm"
      >
        <img
          :src="item.event.featured_image"
          class="h-48 w-full object-cover"
        />
        <div class="p-8 flex flex-col flex-grow">
          <h3 class="text-lg mb-6">{{ item.event.title }}</h3>
          <div class="flex gap-2 mt-auto">
            <NuxtLink
              :to="'/events/' + item.event.id"
              class="flex-1 py-3 bg-slate-900 text-white text-center rounded-xl text-[10px] uppercase"
              >View</NuxtLink
            >
            <button
              @click="handleRemove(item.event.id)"
              class="px-4 py-3 bg-red-50 text-red-500 rounded-xl text-[10px] uppercase"
            >
              Remove
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
