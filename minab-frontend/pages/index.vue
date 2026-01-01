<script setup lang="ts">
import { ref, computed, onMounted } from "vue";
import { searchEventQuery } from "~/utils/constants/queries/events";
import { useForm } from "vee-validate";
import * as zod from "zod";
import { toTypedSchema } from "@vee-validate/zod";

const { data, refresh } = await useAsyncQuery<any>(searchEventQuery, {
  take: 3,
});
const featured = computed(() => data.value?.events || []);

const { handleSubmit, errors, defineField, resetForm } = useForm({
  validationSchema: toTypedSchema(
    zod.object({
      name: zod.string().min(2, "Name is required"),
      email: zod.string().email("Invalid email address"),
      msg: zod.string().min(10, "Message must be at least 10 characters"),
    })
  ),
});

const [name] = defineField("name");
const [email] = defineField("email");
const [msg] = defineField("msg");

const onSend = handleSubmit((values) => {
  console.log("Contact Form Submitted:", values);
  alert(`Thank you ${values.name}! Your message has been sent to Minab Tech.`);
  resetForm(); // Clears the form after sending
});

onMounted(() => refresh());

const heroBg =
  "https://images.unsplash.com/photo-1492684223066-81342ee5ff30?q=80&w=1920";
const fallback =
  "https://images.unsplash.com/photo-1501281668745-f7f57925c3b4?auto=format&fit=crop&q=80&w=1000";
</script>

<template>
  <div class="bg-white font-sans font-normal selection:bg-indigo-100">
    <!-- HERO SECTION -->
    <section
      class="relative h-[65vh] flex items-center justify-center bg-[#0f172a] text-center overflow-hidden"
    >
      <img
        :src="heroBg"
        class="absolute inset-0 w-full h-full object-cover opacity-30"
      />
      <div class="relative z-10 px-6 max-w-4xl text-white">
        <h1
          class="text-4xl md:text-7xl font-normal mb-6 tracking-tighter uppercase leading-tight"
        >
          Unforgettable Events, <br /><span class="text-indigo-400"
            >Just a Click Away.</span
          >
        </h1>
        <div class="flex justify-center gap-4 mt-10">
          <NuxtLink
            to="/events"
            class="px-8 py-3 bg-indigo-600 text-white rounded-xl text-xs uppercase tracking-widest shadow-lg hover:bg-indigo-700 transition-all font-normal"
          >
            Browse Events
          </NuxtLink>
          <NuxtLink
            to="/events/create"
            class="px-8 py-3 border border-white/20 text-white rounded-xl text-xs uppercase tracking-widest hover:bg-white/10 transition-all font-normal"
          >
            Host an Event
          </NuxtLink>
        </div>
      </div>
    </section>

    <section
      class="py-20 px-6 border-b border-slate-50 text-center grid grid-cols-1 md:grid-cols-3 gap-16 max-w-7xl mx-auto"
    >
      <div>
        <div
          class="text-3xl mb-6 bg-blue-50 w-16 h-16 rounded-3xl flex items-center justify-center mx-auto text-blue-500"
        >
          🔍
        </div>
        <h3 class="text-lg font-normal uppercase tracking-widest mb-2">
          Discovery
        </h3>
        <p class="text-slate-400 text-sm font-normal">
          Smart map-based search engine.
        </p>
      </div>
      <div>
        <div
          class="text-3xl mb-6 bg-indigo-50 w-16 h-16 rounded-3xl flex items-center justify-center mx-auto text-indigo-500"
        >
          👤
        </div>
        <h3 class="text-lg font-normal uppercase tracking-widest mb-2">
          Engagement
        </h3>
        <p class="text-slate-400 text-sm font-normal">
          Save bookmarks and manage your tickets.
        </p>
      </div>
      <div>
        <div
          class="text-3xl mb-6 bg-purple-50 w-16 h-16 rounded-3xl flex items-center justify-center mx-auto text-purple-500"
        >
          📊
        </div>
        <h3 class="text-lg font-normal uppercase tracking-widest mb-2">
          Management
        </h3>
        <p class="text-slate-400 text-sm font-normal">
          Host and track events with ease.
        </p>
      </div>
    </section>

    <section class="py-24 px-6 text-center">
      <h2
        class="text-2xl uppercase tracking-[0.3em] text-slate-900 mb-16 font-normal"
      >
        Featured Events
      </h2>
      <div
        class="max-w-7xl mx-auto grid grid-cols-1 md:grid-cols-3 gap-10 text-left"
      >
        <div
          v-for="event in featured"
          :key="event.id"
          class="group bg-white rounded-[2.5rem] border overflow-hidden shadow-sm hover:shadow-xl transition-all h-full flex flex-col"
        >
          <img
            :src="event.featured_image || fallback"
            class="h-56 w-full object-cover group-hover:scale-105 transition-transform duration-1000"
          />
          <div class="p-8 flex flex-col flex-grow">
            <h3
              class="text-lg font-normal text-slate-900 mb-6 capitalize leading-tight"
            >
              {{ event.title }}
            </h3>
            <!-- FIXED LINK -->
            <NuxtLink
              :to="'/events/' + event.id"
              class="mt-auto block w-full py-4 bg-slate-900 text-white text-center rounded-2xl text-[10px] uppercase tracking-widest hover:bg-indigo-600 transition-all"
            >
              View Details →
            </NuxtLink>
          </div>
        </div>
      </div>
    </section>

    <section class="py-24 px-6 bg-slate-50 border-t">
      <div class="max-w-4xl mx-auto text-center font-normal">
        <h2
          class="text-3xl uppercase mb-10 tracking-tighter text-slate-900 font-normal"
        >
          Contact Minab Tech
        </h2>
        <form
          @submit.prevent="onSend"
          class="bg-white p-12 md:p-20 rounded-[4rem] shadow-2xl border space-y-6 text-left"
        >
          <div class="grid grid-cols-1 md:grid-cols-2 gap-8">
            <div class="space-y-1">
              <label
                class="text-[10px] uppercase tracking-widest text-slate-400 ml-4"
                >Full Name</label
              >
              <input
                v-model="name"
                placeholder="Your name"
                class="w-full p-5 bg-slate-50 border-none rounded-2xl outline-none font-normal text-sm focus:ring-1 focus:ring-indigo-400"
              />
              <p class="text-red-400 text-[10px] ml-4 mt-1">
                {{ errors.name }}
              </p>
            </div>
            <div class="space-y-1">
              <label
                class="text-[10px] uppercase tracking-widest text-slate-400 ml-4"
                >Email Address</label
              >
              <input
                v-model="email"
                placeholder="Email"
                class="w-full p-5 bg-slate-50 border-none rounded-2xl outline-none font-normal text-sm focus:ring-1 focus:ring-indigo-400"
              />
              <p class="text-red-400 text-[10px] ml-4 mt-1">
                {{ errors.email }}
              </p>
            </div>
          </div>
          <div class="space-y-1">
            <label
              class="text-[10px] uppercase tracking-widest text-slate-400 ml-4"
              >Message</label
            >
            <textarea
              v-model="msg"
              rows="4"
              placeholder="How can we help?"
              class="w-full p-5 bg-slate-50 border-none rounded-2xl outline-none h-32 font-normal text-sm focus:ring-1 focus:ring-indigo-400"
            ></textarea>
            <p class="text-red-400 text-[10px] ml-4 mt-1">{{ errors.msg }}</p>
          </div>
          <button
            type="submit"
            class="w-full py-5 bg-indigo-600 text-white rounded-[2rem] text-xs uppercase tracking-widest font-normal hover:bg-indigo-700 shadow-xl transition-all active:scale-95"
          >
            Send Message
          </button>
        </form>
        <p
          class="mt-12 text-slate-300 text-[10px] tracking-widest uppercase font-normal"
        >
          support@minabevents.com
        </p>
      </div>
    </section>
  </div>
</template>

<style scoped>
* {
  font-weight: 400 !important;
}
</style>
