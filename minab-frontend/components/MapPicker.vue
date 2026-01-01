<script setup lang="ts">
const props = defineProps(["modelValue"]);
const emit = defineEmits(["update:modelValue", "addressSelected"]);

const center = ref([9.01, 38.75]); // Addis Ababa
const markerPos = ref([9.01, 38.75]);
const zoom = ref(12);

const handleMapClick = async (e: any) => {
  const { lat, lng } = e.latlng;
  markerPos.value = [lat, lng];

  // Emit coordinates to the parent form
  emit("update:modelValue", { lat, lng });

  // Professional Touch: Auto-fetch address using OpenStreetMap Nominatim
  try {
    const res = await fetch(
      `https://nominatim.openstreetmap.org/reverse?format=json&lat=${lat}&lon=${lng}`
    );
    const data = await res.json();
    emit("addressSelected", data.display_name);
  } catch (error) {
    console.error("Geocoding failed", error);
  }
};
</script>

<template>
  <div
    class="rounded-3xl overflow-hidden border-2 border-slate-100 shadow-inner h-80 w-full relative"
  >
    <LMap
      ref="map"
      :zoom="zoom"
      :center="center"
      :use-global-leaflet="false"
      @click="handleMapClick"
    >
      <LTileLayer
        url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png"
        layer-type="base"
        name="OpenStreetMap"
      />
      <LMarker :lat-lng="markerPos" />
    </LMap>
    <div
      class="absolute bottom-4 left-4 z-[1000] bg-white/90 backdrop-blur-md px-4 py-2 rounded-xl text-[10px] font-bold uppercase tracking-widest text-slate-500 shadow-sm border border-slate-200"
    >
      Click on the map to set venue location
    </div>
  </div>
</template>
