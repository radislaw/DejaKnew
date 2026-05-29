<script setup lang="ts">
defineProps<{
  flipped?: boolean
}>()

defineEmits<{
  flip: []
}>()
</script>

<template>
  <div
    class="flip-card"
    :class="{ 'is-flipped': flipped }"
    @click="$emit('flip')"
  >
    <div class="flip-card__inner">
      <div class="flip-card__face flip-card__face--front">
        <slot name="front" />
      </div>
      <div class="flip-card__face flip-card__face--back">
        <slot name="back" />
      </div>
    </div>
  </div>
</template>

<style scoped>
.flip-card {
  display: block;
  transform-style: preserve-3d;
}

.flip-card__inner {
  position: relative;
  width: 100%;
  height: 100%;
  transform-style: preserve-3d;
  transition: transform 1.4s cubic-bezier(.5, .05, .2, 1);
}

.flip-card.is-flipped .flip-card__inner {
  transform: rotateY(180deg);
}

.flip-card__face {
  position: absolute;
  inset: 0;
  backface-visibility: hidden;
  -webkit-backface-visibility: hidden;
}

.flip-card__face--back {
  transform: rotateY(180deg);
}
</style>
