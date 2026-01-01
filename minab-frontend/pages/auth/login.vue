<script setup lang="ts">
import { ref } from "vue";
import { useForm } from "vee-validate";
import { required, isValidEmail } from "../../utils/helpers/validation";
import { loginMutation } from "../../utils/constants/queries/auth";
import ErrorIcon from "../../components/icons/Error.vue";
import CloseIcon from "../../components/icons/Close.vue";
import LoadingIcon from "../../components/icons/Loading.vue";

definePageMeta({
  layout: "auth",
});

// 1. Setup Form Validation
const { defineField, handleSubmit, errors } = useForm({
  validationSchema: {
    email: isValidEmail,
    password: required,
  },
});

const [email, emailProps] = defineField("email");
const [password, passwordProps] = defineField("password");

// 2. Setup States
const rememberMe = ref(false);
const errorMessage = ref("");
const isLoading = ref(false);
const isError = ref(false);

// 3. Use Apollo Mutation (Technical Requirement)
const { mutate: loginUser } = useMutation(loginMutation);

const onSubmit = handleSubmit(async () => {
  isLoading.value = true;
  try {
    const res: any = await loginUser({
      email: email.value, // Changed from login_text to email
      password: password.value,
      remember_me: rememberMe.value,
    });

    if (res?.data?.login?.token) {
      const token = useCookie("token", { maxAge: 60 * 60 * 24 * 7 });
      token.value = res.data.login.token;
      navigateTo("/");
    }
  } catch (err: any) {
    errorMessage.value = err.message;
    isError.value = true;
  } finally {
    isLoading.value = false;
  }
});
const toggle = () => {
  isError.value = false;
};
</script>

<template>
  <div class="min-h-screen flex items-center justify-center bg-slate-50 px-4">
    <!-- Error Toast -->
    <Transition name="fade">
      <div
        v-if="isError"
        class="fixed top-5 right-5 flex items-center w-full max-w-xs p-4 text-slate-500 bg-white rounded-2xl shadow-xl border border-red-100 z-50"
        role="alert"
      >
        <div
          class="inline-flex items-center justify-center flex-shrink-0 w-10 h-10 text-red-500 bg-red-50 rounded-xl"
        >
          <ErrorIcon class="w-6 h-6" />
        </div>
        <div class="ms-3 text-sm font-bold text-slate-700">
          {{ errorMessage }}
        </div>
        <button
          @click="toggle"
          class="ms-auto bg-transparent text-slate-400 hover:text-slate-900 p-1.5"
        >
          <CloseIcon class="w-5 h-5" />
        </button>
      </div>
    </Transition>

    <div class="w-full max-w-md">
      <!-- Card Container -->
      <div
        class="bg-white p-10 rounded-[2.5rem] shadow-sm border border-slate-200"
      >
        <div class="flex flex-col items-center mb-8">
          <NuxtLink to="/">
            <img
              src="../../assets/images/logo.png"
              class="h-14 mb-6"
              alt="Minab Logo"
            />
          </NuxtLink>
          <h1 class="text-3xl font-black text-slate-900 tracking-tight">
            Welcome Back
          </h1>
          <p class="text-slate-500 text-sm mt-2">
            Login to manage your events and tickets
          </p>
        </div>

        <form class="space-y-5" @submit.prevent="onSubmit">
          <!-- Email Field -->
          <div>
            <label
              for="email"
              class="block mb-2 text-xs font-black uppercase tracking-widest text-slate-400"
              >Email Address</label
            >
            <input
              v-model="email"
              v-bind="emailProps"
              type="email"
              id="email"
              class="w-full p-4 bg-slate-50 border-none rounded-2xl focus:ring-2 focus:ring-indigo-500 outline-none transition-all placeholder:text-slate-300"
              placeholder="name@company.com"
            />
            <span class="text-[10px] text-red-500 font-bold mt-1 block">{{
              errors.email
            }}</span>
          </div>

          <!-- Password Field -->
          <div>
            <label
              for="password"
              class="block mb-2 text-xs font-black uppercase tracking-widest text-slate-400"
              >Password</label
            >
            <input
              v-model="password"
              v-bind="passwordProps"
              type="password"
              id="password"
              placeholder="••••••••"
              class="w-full p-4 bg-slate-50 border-none rounded-2xl focus:ring-2 focus:ring-indigo-500 outline-none transition-all placeholder:text-slate-300"
            />
            <span class="text-[10px] text-red-500 font-bold mt-1 block">{{
              errors.password
            }}</span>
          </div>

          <!-- Remember Me -->
          <div class="flex items-center justify-between">
            <div class="flex items-center">
              <input
                id="remember"
                type="checkbox"
                v-model="rememberMe"
                class="w-4 h-4 text-indigo-600 border-slate-300 rounded focus:ring-indigo-500"
              />
              <label
                for="remember"
                class="ms-2 text-sm font-bold text-slate-600"
                >Remember me</label
              >
            </div>
            <a
              href="#"
              class="text-sm font-bold text-indigo-600 hover:underline"
              >Forgot password?</a
            >
          </div>

          <!-- Submit Button -->
          <button
            type="submit"
            :disabled="isLoading"
            class="w-full py-4 bg-indigo-600 hover:bg-indigo-700 text-white rounded-2xl font-bold text-lg transition-all shadow-lg shadow-indigo-200 disabled:opacity-50 active:scale-95 flex items-center justify-center gap-2"
          >
            <LoadingIcon v-if="isLoading" class="w-5 h-5 animate-spin" />
            {{ isLoading ? "Verifying..." : "Sign In" }}
          </button>

          <!-- Signup Link -->
          <div class="text-center pt-4">
            <p class="text-sm font-bold text-slate-500">
              New to Minab?
              <NuxtLink
                to="/auth/signup"
                class="text-indigo-600 hover:underline ml-1"
                >Create an account</NuxtLink
              >
            </p>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<style scoped>
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.3s ease;
}
.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>
