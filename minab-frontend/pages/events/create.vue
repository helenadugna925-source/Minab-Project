<script setup lang="ts">
import { ref, computed } from "vue";
import { useForm } from "vee-validate";
import * as zod from "zod";
import { toTypedSchema } from "@vee-validate/zod";
import "leaflet/dist/leaflet.css";
import { CREATE_EVENT_MUTATION } from "~/utils/constants/queries/events";
import { uploadFileMutation } from "~/utils/constants/queries/upload";
import { GetCategories } from "~/utils/constants/queries/category";

const token = useCookie("token");

const schema = toTypedSchema(
  zod.object({
    title: zod.string().min(5, "Title too short"),
    description: zod.string().min(10, "Description too short"),
    price: zod.number().min(0),
    date: zod.string().min(1, "Date required"),
    venue_name: zod.string().min(1, "Venue required"),
    address: zod.string().min(1, "Address required"),
    category_id: zod.number(),
  })
);

const { defineField, handleSubmit, errors } = useForm({
  validationSchema: schema,
  initialValues: {
    price: 0,
    category_id: 1,
    title: "",
    description: "",
    date: "",
    venue_name: "",
    address: "",
  },
});

const [title] = defineField("title");
const [description] = defineField("description");
const [price] = defineField("price");
const [date] = defineField("date");
const [venue_name] = defineField("venue_name");
const [address] = defineField("address");
const [category_id] = defineField("category_id");

const { mutate: createEvent } = useMutation(CREATE_EVENT_MUTATION);
const { mutate: uploadFile } = useMutation(uploadFileMutation);
const { data: catData } = await useAsyncQuery<any>(GetCategories);
const categoryList = computed(() => catData.value?.categories || []);

const uploadedUrls = ref<string[]>([]);
const featuredImg = ref("");
const lat = ref(9.0122);
const lng = ref(38.7578);
const isLoading = ref(false);

const onMapClick = (e: any) => {
  lat.value = e.latlng.lat;
  lng.value = e.latlng.lng;
};

const handleUpload = async (e: any) => {
  const files = e.target.files;
  for (let file of files) {
    const base64 = await new Promise((res) => {
      const r = new FileReader();
      r.readAsDataURL(file);
      r.onload = () => res((r.result as string).split(",")[1]);
    });
    const res: any = await uploadFile({ name: file.name, base64 });
    if (res.data?.uploadFile?.url) {
      uploadedUrls.value.push(res.data.uploadFile.url);
      if (!featuredImg.value) featuredImg.value = res.data.uploadFile.url;
    }
  }
};

const onSubmit = handleSubmit(async (values) => {
  if (uploadedUrls.value.length === 0)
    return alert("Please upload at least one image!");
  isLoading.value = true;
  try {
    await createEvent(
      {
        ...values,
        location_lat: parseFloat(lat.value.toString()),
        location_lng: parseFloat(lng.value.toString()),
        tags: ["General"],
        image_urls: uploadedUrls.value,
        featured_image: featuredImg.value || uploadedUrls.value[0], // FIXED
      },
      {
        context: { headers: { Authorization: `Bearer ${token.value}` } },
      }
    );
    navigateTo("/events/my");
  } catch (err: any) {
    alert(err.message);
  } finally {
    isLoading.value = false;
  }
});
</script>

<template>
  <div class="min-h-screen bg-slate-50 py-12 px-6 font-normal">
    <div class="max-w-6xl mx-auto flex justify-between items-center mb-12">
      <h1 class="text-3xl uppercase tracking-tight">Host an Event</h1>
      <NuxtLink to="/events/my" class="px-6 py-2 border rounded-xl text-xs"
        >Cancel</NuxtLink
      >
    </div>

    <form @submit.prevent="onSubmit" class="grid lg:grid-cols-12 gap-10">
      <div class="lg:col-span-8 space-y-10">
        <!-- 01 INFO -->
        <div class="bg-white p-10 rounded-[3rem] border space-y-6 shadow-sm">
          <h2 class="text-[10px] text-indigo-600 uppercase tracking-[0.4em]">
            01. Information
          </h2>
          <input
            v-model="title"
            class="w-full p-4 bg-slate-50 rounded-2xl outline-none"
            placeholder="Event Title"
          />
          <textarea
            v-model="description"
            rows="5"
            class="w-full p-4 bg-slate-50 rounded-2xl outline-none h-40"
            placeholder="Description"
          ></textarea>
          <div class="grid grid-cols-2 gap-4">
            <select
              v-model="category_id"
              class="p-4 bg-slate-50 rounded-2xl outline-none font-normal"
            >
              <option v-for="c in categoryList" :key="c.id" :value="c.id">
                {{ c.name }}
              </option>
            </select>
            <input
              v-model="date"
              type="datetime-local"
              class="p-4 bg-slate-50 rounded-2xl outline-none font-normal"
            />
          </div>
        </div>

        <div class="bg-white p-10 rounded-[3rem] border">
          <h2
            class="text-[10px] text-indigo-600 uppercase tracking-[0.4em] mb-6"
          >
            02. Media
          </h2>
          <input
            type="file"
            multiple
            @change="handleUpload"
            class="mb-6 block w-full text-sm text-slate-500 file:mr-4 file:py-2 file:px-4 file:rounded-full file:border-0 file:text-xs file:font-semibold file:bg-indigo-50 file:text-indigo-700"
          />
          <div class="flex gap-4 flex-wrap">
            <div
              v-for="url in uploadedUrls"
              :key="url"
              @click="featuredImg = url"
              class="relative w-24 h-24 rounded-2xl overflow-hidden border-2 cursor-pointer transition-all"
              :class="
                featuredImg === url ? 'border-indigo-600' : 'border-transparent'
              "
            >
              <img :src="url" class="w-full h-full object-cover" />
              <div
                v-if="featuredImg === url"
                class="absolute top-1 right-1 bg-indigo-600 text-white p-1 rounded-full text-[8px]"
              >
                ★
              </div>
            </div>
          </div>
        </div>
      </div>

      <div class="lg:col-span-4 space-y-6">
        <div
          class="bg-white p-8 rounded-[3rem] border sticky top-24 space-y-6 shadow-sm"
        >
          <h2 class="text-[10px] text-indigo-600 uppercase tracking-[0.4em]">
            03. Venue & Location
          </h2>
          <input
            v-model="venue_name"
            placeholder="Venue Name"
            class="w-full p-4 bg-slate-50 rounded-2xl outline-none"
          />
          <input
            v-model="address"
            placeholder="Full Address"
            class="w-full p-4 bg-slate-50 rounded-2xl outline-none"
          />

          <div
            class="h-44 bg-slate-100 rounded-[2rem] overflow-hidden relative"
          >
            <client-only>
              <LMap
                :zoom="13"
                :center="[lat, lng]"
                :use-global-leaflet="false"
                @click="onMapClick"
              >
                <LTileLayer
                  url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png"
                />
                <LMarker :lat-lng="[lat, lng]" />
              </LMap>
            </client-only>
          </div>

          <input
            v-model="price"
            type="number"
            placeholder="Price (ETB)"
            class="w-full p-4 bg-slate-50 rounded-2xl font-normal"
          />
          <button
            type="submit"
            :disabled="isLoading"
            class="w-full py-5 bg-indigo-600 text-white rounded-[2rem] text-sm uppercase tracking-widest font-normal hover:bg-indigo-700 transition-all"
          >
            Publish Event
          </button>
        </div>
      </div>
    </form>
  </div>
</template>
