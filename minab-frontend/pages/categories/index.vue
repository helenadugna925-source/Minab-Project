<script setup lang="ts">
import { useForm } from "vee-validate";
import { required } from "~/utils/helpers/validation";

definePageMeta({
  layout: "auth",
});

const config = useRuntimeConfig();
const apiBase = config.public.API_BASE;
const token = useCookie("token");
const categories = ref<{ id: number; name: string; image?: string }[]>([]);
const isLoading = ref(false);
const isError = ref(false);
const message = ref("");

const fetchCategories = async () => {
  isLoading.value = true;
  try {
    const res: any = await $fetch(`${config.public.GQL_HOST}`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        Authorization: token.value ? `Bearer ${token.value}` : "",
      },
      body: {
        query: `
          query {
            categories {
              id
              name
              image
            }
          }
        `,
      },
    });
    categories.value = res.data?.categories || [];
  } catch (err) {
    console.log(err);
    isError.value = true;
    message.value = "Failed to load categories";
  } finally {
    isLoading.value = false;
  }
};

await fetchCategories();

const { defineField, handleSubmit, errors, resetForm } = useForm({
  validationSchema: {
    name: required,
  },
});
const [name, nameProps] = defineField("name");

const createCategory = handleSubmit(async () => {
  try {
    await $fetch(`${config.public.GQL_HOST}`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        Authorization: token.value ? `Bearer ${token.value}` : "",
      },
      body: {
        query: `
          mutation($name: String!) {
            insert_categories_one(object: {name: $name}) { id }
          }
        `,
        variables: { name: name.value },
      },
    });
    message.value = "Category created";
    await fetchCategories();
    resetForm();
  } catch (err) {
    console.log(err);
    isError.value = true;
    message.value = "Failed to create category";
  }
});

const deleteCategory = async (id: number) => {
  try {
    await $fetch(`${config.public.GQL_HOST}`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        Authorization: token.value ? `Bearer ${token.value}` : "",
      },
      body: {
        query: `
          mutation($id: Int!) {
            delete_categories_by_pk(id: $id) { id }
          }
        `,
        variables: { id },
      },
    });
    message.value = "Category deleted";
    await fetchCategories();
  } catch (err) {
    console.log(err);
    isError.value = true;
    message.value = "Failed to delete category";
  }
};
</script>

<template>
  <div class="max-w-3xl mx-auto mt-8">
    <h2 class="text-2xl font-bold mb-4">Categories</h2>

    <form class="mb-6" @submit.prevent="createCategory">
      <label class="block mb-2 text-sm font-medium text-text-default"
        >Name</label
      >
      <input
        type="text"
        v-model="name"
        v-bind="nameProps"
        class="bg-bg-default border border-gray-300 text-text-default text-sm rounded-lg focus:ring-accent focus:border-accent block w-full p-2.5"
        placeholder="Category name"
      />
      <span class="text-sm text-red-600">{{ errors.name }}</span>
      <button
        type="submit"
        class="mt-3 text-white bg-primary hover:bg-accent focus:ring-4 focus:outline-none focus:ring-accent font-medium rounded-lg text-sm px-5 py-2.5 text-center"
      >
        Add Category
      </button>
    </form>

    <div
      v-if="message"
      class="mb-4 text-sm"
      :class="isError ? 'text-red-600' : 'text-green-600'"
    >
      {{ message }}
    </div>

    <div v-if="isLoading" class="text-muted">Loading...</div>
    <div v-else class="space-y-2">
      <div
        v-for="cat in categories"
        :key="cat.id"
        class="flex items-center justify-between p-3 border rounded"
      >
        <div>
          <p class="font-medium">{{ cat.name }}</p>
        </div>
        <button
          @click="deleteCategory(cat.id)"
          class="text-red-600 border border-red-600 px-3 py-1 rounded hover:bg-red-50"
        >
          Delete
        </button>
      </div>
    </div>
  </div>
  <Footer />
</template>
