<script setup lang="ts">
import { useColorMode } from '@vueuse/core'
import { ref, watch } from 'vue'
import type { Ref } from 'vue'
import { Box, Check, Moon, Search, Sun } from 'lucide-vue-next'
import { Button } from '@/components/ui/button'
import {
  Select,
  SelectItem,
  SelectTrigger,
  SelectContent,
  SelectGroup,
  SelectValue,
} from '@/components/ui/select'
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from '@/components/ui/card'
import {
  TagsInput,
  TagsInputInput,
  TagsInputItem,
  TagsInputItemDelete,
  TagsInputItemText,
} from '@/components/ui/tags-input'
import { Skeleton } from '@/components/ui/skeleton'
import {
  Stepper,
  StepperDescription,
  StepperItem,
  StepperSeparator,
  StepperTrigger,
  StepperTitle,
} from '@/components/ui/stepper'

const apiCorreios = import.meta.env.VITE_API_URL

const mode = useColorMode()
const buscaObjetos: Ref<string[]> = ref([])
const objetoSelecionado: Ref<Objeto | undefined> = ref(undefined)
const objetos: Ref<Objeto[]> = ref([])

watch(buscaObjetos, async (novoBuscaObjetos) => {
  console.log(apiCorreios)
  const objetosBuscadosPromessa = await fetch(apiCorreios, {
    method: 'POST',
    headers: {
      Accept: 'application/json',
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({
      objetos: novoBuscaObjetos,
    }),
  })

  const objetosBuscados: Objeto[] = await objetosBuscadosPromessa.json()
  objetos.value = objetosBuscados
})

watch(objetos, async (novosObjetos) => {
  if (novosObjetos.length === 0) {
    objetoSelecionado.value = undefined
  }
})
</script>

<template>
  <div class="flex flex-col items-center mt-8">
    <header class="relative mb-8 w-120">
      <Box class="block mt-0 mx-auto size-[125px] text-muted-foreground" />
      <Button
        variant="ghost"
        @click="mode === 'dark' ? (mode = 'light') : (mode = 'dark')"
        class="absolute top-0 right-0 text-muted-foreground dark:text-white"
      >
        <Sun v-if="mode === 'dark'" class="size-5" />
        <Moon v-else class="size-5" />
      </Button>
    </header>

    <main class="flex flex-col gap-2 w-120">
      <div class="flex gap-2">
        <div class="relative w-full items-center">
          <TagsInput
            v-model="buscaObjetos"
            id="search"
            type="text"
            :addOnTab="true"
            :addOnPaste="true"
            :delimiter="RegExp('[., ;]')"
            placeholder="Insira os códigos de rastreamento aqui"
            class="h-full pl-10"
          >
            <TagsInputItem v-for="objeto in buscaObjetos" :key="objeto" :value="objeto">
              <TagsInputItemText />
              <TagsInputItemDelete />
            </TagsInputItem>

            <TagsInputInput
              :autoFocus="true"
              placeholder="Insira os códigos de rastreamento aqui"
            />
          </TagsInput>
          <span class="absolute start-0 inset-y-0 flex items-center justify-center px-2">
            <Search class="size-5 text-muted-foreground dark:text-white" />
          </span>
        </div>
        <Button type="submit">Buscar</Button>
      </div>

      <Select v-model="objetoSelecionado">
        <SelectTrigger class="w-full" :disabled="objetos === undefined || objetos.length === 0">
          <SelectValue placeholder="Selecione um objeto" />
        </SelectTrigger>
        <SelectContent>
          <SelectGroup>
            <SelectItem v-for="objeto in objetos" :key="objeto.codigo" :value="objeto">
              {{ objeto.codigo }}
            </SelectItem>
          </SelectGroup>
        </SelectContent>
      </Select>

      <ul class="flex flex-col gap-2">
        <li v-if="!objetoSelecionado">
          <Card>
            <CardHeader>
              <CardTitle><Skeleton class="h-4 w-28" /></CardTitle>
              <CardDescription><Skeleton class="h-5 w-20" /></CardDescription>
            </CardHeader>
            <CardContent>
              <Stepper
                orientation="vertical"
                class="mx-auto flex w-full flex-col justify-start gap-10"
              >
                <StepperItem class="relative flex w-full items-start gap-6" :step="1">
                  <StepperTrigger as-child>
                    <Skeleton class="size-9 z-10 rounded-full shrink-0" />
                  </StepperTrigger>

                  <div class="flex flex-col gap-1 w-full">
                    <div class="flex items-center justify-between">
                      <StepperTitle>
                        <Skeleton class="h-5 lg:h-6 w-10" />
                      </StepperTitle>
                      <span>
                        <Skeleton class="h-5 w-30" />
                      </span>
                    </div>
                    <StepperDescription>
                      <Skeleton class="h-4 lg:h-5 w-40" />
                    </StepperDescription>
                  </div>
                </StepperItem>
              </Stepper>
            </CardContent>
            <CardFooter>
              <div class="h-4 lg:h-5 w-full">
                <Skeleton class="float-right h-full w-16" />
              </div>
            </CardFooter>
          </Card>
        </li>
        <li v-else>
          <Card>
            <CardHeader>
              <CardTitle>
                {{ objetoSelecionado.codigo }}
              </CardTitle>
              <CardDescription>{{ objetoSelecionado.mensagem }}</CardDescription>
            </CardHeader>
            <CardContent>
              <Stepper
                orientation="vertical"
                class="mx-auto flex w-full flex-col justify-start gap-10"
              >
                <StepperItem
                  v-for="evento in objetoSelecionado.eventos"
                  :key="evento.statusEvento"
                  class="relative flex w-full items-start gap-6"
                  :step="parseInt(evento.statusEvento)"
                >
                  <StepperSeparator
                    v-if="
                      evento.statusEvento !==
                      objetoSelecionado.eventos[objetoSelecionado.eventos.length - 1].statusEvento
                    "
                    class="absolute left-[18px] top-[38px] block h-[105%] w-0.5 shrink-0 rounded-full bg-muted group-data-[state=completed]:bg-primary"
                  />

                  <StepperTrigger as-child>
                    <Button variant="default" size="icon" class="z-10 rounded-full shrink-0">
                      <Check class="size-5" />
                    </Button>
                  </StepperTrigger>

                  <div class="flex flex-col gap-1 w-full">
                    <div class="flex items-center justify-between">
                      <StepperTitle
                        class="text-sm font-semibold transition lg:text-base text-primary"
                      >
                        {{ evento.tipoEvento }}
                      </StepperTitle>
                      <span class="text-muted-foreground text-sm">
                        {{ evento.dataCriacao }}
                      </span>
                    </div>
                    <StepperDescription
                      class="sr-only text-xs text-muted-foreground transition md:not-sr-only lg:text-sm text-primary"
                    >
                      {{ evento.descricaoEvento }}
                    </StepperDescription>
                  </div>
                </StepperItem>
              </Stepper>
            </CardContent>
            <CardFooter class="text-muted-foreground text-xs lg:text-sm text-end">
              <span v-if="objetoSelecionado.eventos[0].paisRemetente === 'Brasil'" class="w-full">
                NACIONAL
              </span>
              <span v-else class="w-full"> INTERNACIONAL </span>
            </CardFooter>
          </Card>
        </li>
      </ul>
    </main>
  </div>
</template>
