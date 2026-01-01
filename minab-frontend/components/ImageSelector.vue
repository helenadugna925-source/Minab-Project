<script setup lang="ts">
const props = defineProps(["images", "featuredImage"]);
const emit = defineEmits(["upload", "setFeatured", "remove"]);

const onFileChange = (e: any) => {
  const files = Array.from(e.target.files);
  emit("upload", files);
};
</script>

<template>
  <div class="space-y-6">
    <!-- Upload Button -->
    <label
      class="flex flex-col items-center justify-center w-full h-40 border-2 border-dashed border-slate-200 rounded-[2rem] cursor-pointer hover:bg-slate-50 transition-all group"
    >
      <div class="flex flex-col items-center justify-center pt-5 pb-6">
        <svg
          class="w-10 h-10 mb-3 text-slate-400 group-hover:text-indigo-500 transition-colors"
          fill="none"
          stroke="currentColor"
          viewBox="0 0 24 24"
        >
          <path
            stroke-linecap="round"
            stroke-linejoin="round"
            stroke-width="2"
            d="M12 4v16m8-8H4"
          ></path>
        </svg>
        <p class="text-sm text-slate-500 font-medium">
          Click to upload event photos
        </p>
        <p class="text-xs text-slate-400 mt-1">PNG, JPG up to 10MB</p>
      </div>
      <input type="file" multiple class="hidden" @change="onFileChange" />
    </label>

    <!-- Image Preview Grid -->
    <div v-if="images.length > 0" class="grid grid-cols-2 md:grid-cols-4 gap-4">
      <div
        v-for="(url, index) in images"
        :key="index"
        class="relative group aspect-square rounded-2xl overflow-hidden border-2 transition-all"
        :class="
          featuredImage === url
            ? 'border-indigo-600 ring-4 ring-indigo-50'
            : 'border-transparent'
        "
      >
        <img :src="url" class="w-full h-full object-cover" />

        <!-- Overlay -->
        <div
          class="absolute inset-0 bg-black/40 opacity-0 group-hover:opacity-100 transition-opacity flex flex-col items-center justify-center gap-2"
        >
          <button
            @click="emit('setFeatured', url)"
            class="px-3 py-1.5 bg-white text-indigo-600 text-[10px] font-bold rounded-lg uppercase"
          >
            {{ featuredImage === url ? "★ Featured" : "Set Featured" }}
          </button>
          <button
            @click="emit('remove', index)"
            class="text-white text-[10px] font-bold hover:text-red-400"
          >
            Remove
          </button>
        </div>

        <!-- Featured Badge -->
        <div
          v-if="featuredImage === url"
          class="absolute top-2 left-2 bg-indigo-600 text-white text-[8px] font-black px-2 py-1 rounded-md uppercase"
        >
          Thumbnail
        </div>
      </div>
    </div>
  </div>
</template>
