<script setup lang="ts">
import { ref, watch } from "vue";
import gql from "graphql-tag";

const MAP_QUERY = gql`
  query GetEventsOnMap(
    $minLat: Float!
    $maxLat: Float!
    $minLng: Float!
    $maxLng: Float!
  ) {
    get_events_on_map(
      args: {
        min_lat: $minLat
        max_lat: $maxLat
        min_lng: $minLng
        max_lng: $maxLng
      }
    ) {
      id
      title
      featured_image
      price
      location_lat
      location_lng
      venue_name
      event_date
    }
  }
`;

const mapBounds = ref({ minLat: 8.5, maxLat: 9.5, minLng: 38.5, maxLng: 39.5 });
const hoveredEventId = ref<number | null>(null);

const { result, loading, refetch } = useQuery(MAP_QUERY, mapBounds.value);

const onMapMoveEnd = (event: any) => {
  const bounds = event.target.getBounds();
  mapBounds.value = {
    minLat: bounds.getSouth(),
    maxLat: bounds.getNorth(),
    minLng: bounds.getWest(),
    maxLng: bounds.getEast(),
  };
  refetch(mapBounds.value);
};

const formatDate = (dateStr: string) =>
  new Date(dateStr).toLocaleDateString("en-US", {
    month: "short",
    day: "numeric",
  });
</script>

<template>
  <div class="flex h-[calc(100vh-64px)] overflow-hidden">
    <!-- Left Side: Interactive Map -->
    <div class="w-2/3 h-full relative border-r">
      <LMap
        style="height: 100%"
        :zoom="12"
        :center="[9.01, 38.75]"
        :use-global-leaflet="false"
        @moveend="onMapMoveEnd"
      >
        <LTileLayer url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png" />

        <!-- Professional Markers -->
        <LMarker
          v-for="event in result?.get_events_on_map"
          :key="event.id"
          :lat-lng="[event.location_lat, event.location_lng]"
        >
          <LPopup>
            <div class="p-1 w-40">
              <img
                :src="event.featured_image"
                class="w-full h-20 object-cover rounded mb-2"
              />
              <h4 class="font-bold text-sm leading-tight">{{ event.title }}</h4>
              <p class="text-xs text-brand-orange font-bold mt-1">
                {{ event.price === 0 ? "FREE" : event.price + " ETB" }}
              </p>
            </div>
          </LPopup>
        </LMarker>
      </LMap>

      <div
        v-if="loading"
        class="absolute top-4 left-1/2 -translate-x-1/2 z-[1000] bg-white px-4 py-2 rounded-full shadow-lg text-sm font-medium flex items-center"
      >
        <span class="animate-spin mr-2">⏳</span> Searching area...
      </div>
    </div>

    <!-- Right Side: Scrollable Results (Professional Grid) -->
    <div class="w-1/3 h-full overflow-y-auto bg-gray-50 p-4">
      <div class="flex justify-between items-center mb-6">
        <h2 class="text-xl font-bold text-brand-dark">Events in this area</h2>
        <span class="text-sm text-gray-500"
          >{{ result?.get_events_on_map?.length || 0 }} results</span
        >
      </div>

      <div class="space-y-4">
        <div
          v-for="event in result?.get_events_on_map"
          :key="event.id"
          @mouseenter="hoveredEventId = event.id"
          @mouseleave="hoveredEventId = null"
          class="bg-white rounded-xl shadow-sm border overflow-hidden hover:shadow-md transition-all cursor-pointer flex"
        >
          <img :src="event.featured_image" class="w-32 h-32 object-cover" />
          <div class="p-3 flex flex-col justify-between flex-1">
            <div>
              <p
                class="text-[10px] font-bold text-brand-orange uppercase tracking-wider"
              >
                {{ formatDate(event.event_date) }}
              </p>
              <h3
                class="font-bold text-brand-dark line-clamp-2 leading-tight mt-1"
              >
                {{ event.title }}
              </h3>
              <p class="text-xs text-gray-500 mt-1">{{ event.venue_name }}</p>
            </div>
            <div class="flex justify-between items-center mt-2">
              <span class="text-sm font-bold text-gray-900">{{
                event.price === 0 ? "Free" : event.price + " ETB"
              }}</span>
              <button class="text-brand-blue text-xs font-bold hover:underline">
                Details →
              </button>
            </div>
          </div>
        </div>

        <!-- Empty State -->
        <div
          v-if="!loading && result?.get_events_on_map?.length === 0"
          class="text-center py-20"
        >
          <div class="text-4xl mb-4">📍</div>
          <p class="text-gray-500">
            No events found in this area.<br />Try moving the map!
          </p>
        </div>
      </div>
    </div>
  </div>
</template>
