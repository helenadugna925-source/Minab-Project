<script setup lang="ts">
import { useForm } from "vee-validate";
import * as zod from "zod";
import { toTypedSchema } from "@vee-validate/zod";

// 1. Validation Schema
const contactSchema = toTypedSchema(
  zod.object({
    name: zod.string().min(2, "Name is too short"),
    email: zod.string().email("Invalid email"),
    message: zod.string().min(10, "Message must be at least 10 characters"),
  })
);

const { handleSubmit, errors, defineField, resetForm } = useForm({
  validationSchema: contactSchema,
});

const [name] = defineField("name");
const [email] = defineField("email");
const [message] = defineField("message");

// 2. Submission Logic
const onSendContact = handleSubmit(async (values) => {
  console.log("Contact Form Data:", values);
  // Simulating an API call
  alert("Message Sent! Thank you for contacting Minab Events.");
  resetForm();
});
</script>

<template>
  <div
    class="max-w-4xl mx-auto bg-white p-12 rounded-[3rem] shadow-2xl border border-slate-50 my-20"
  >
    <form @submit="onSendContact" class="space-y-6">
      <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
        <div>
          <label
            class="block text-xs font-normal text-slate-400 uppercase tracking-widest mb-2"
            >Name</label
          >
          <input
            v-model="name"
            class="w-full p-4 bg-slate-50 border-none rounded-2xl outline-none focus:ring-2 focus:ring-indigo-500"
            placeholder="Your name"
          />
          <p class="text-red-400 text-[10px] mt-1">{{ errors.name }}</p>
        </div>
        <div>
          <label
            class="block text-xs font-normal text-slate-400 uppercase tracking-widest mb-2"
            >Email</label
          >
          <input
            v-model="email"
            class="w-full p-4 bg-slate-50 border-none rounded-2xl outline-none focus:ring-2 focus:ring-indigo-500"
            placeholder="Your email"
          />
          <p class="text-red-400 text-[10px] mt-1">{{ errors.email }}</p>
        </div>
      </div>
      <div>
        <label
          class="block text-xs font-normal text-slate-400 uppercase tracking-widest mb-2"
          >Message</label
        >
        <textarea
          v-model="message"
          rows="4"
          class="w-full p-4 bg-slate-50 border-none rounded-2xl outline-none focus:ring-2 focus:ring-indigo-500"
          placeholder="How can we help?"
        ></textarea>
        <p class="text-red-400 text-[10px] mt-1">{{ errors.message }}</p>
      </div>

      <button
        type="submit"
        class="w-full py-5 bg-sky-500 hover:bg-sky-600 text-white rounded-2xl font-normal text-sm uppercase tracking-[0.2em] shadow-lg shadow-sky-100 transition-all active:scale-95"
      >
        Send Message
      </button>
    </form>
  </div>
</template>
