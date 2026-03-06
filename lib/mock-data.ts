import { Profile, Match, Message } from "./types"

export const mockProfiles: Profile[] = [
  {
    id: "1",
    name: "Sofía",
    age: 26,
    bio: "Amante del café, los viajes y los buenos libros. Buscando a alguien con quien compartir atardeceres.",
    images: [
      "https://images.unsplash.com/photo-1494790108377-be9c29b29330?w=800&q=80",
      "https://images.unsplash.com/photo-1524504388940-b1c1722653e1?w=800&q=80",
    ],
    distance: 3,
    occupation: "Diseñadora UX",
    interests: ["Viajes", "Fotografía", "Yoga", "Café"],
  },
  {
    id: "2",
    name: "Valentina",
    age: 24,
    bio: "Ingeniera de día, chef experimental de noche. Me encanta probar nuevos restaurantes.",
    images: [
      "https://images.unsplash.com/photo-1517841905240-472988babdf9?w=800&q=80",
      "https://images.unsplash.com/photo-1534528741775-53994a69daeb?w=800&q=80",
    ],
    distance: 5,
    occupation: "Ingeniera de Software",
    interests: ["Cocina", "Tecnología", "Running", "Música"],
  },
  {
    id: "3",
    name: "Camila",
    age: 28,
    bio: "Médica apasionada por la naturaleza. Los fines de semana me encontrarás en la montaña.",
    images: [
      "https://images.unsplash.com/photo-1488426862026-3ee34a7d66df?w=800&q=80",
      "https://images.unsplash.com/photo-1507003211169-0a1dd7228f2d?w=800&q=80",
    ],
    distance: 8,
    occupation: "Médica",
    interests: ["Senderismo", "Mascotas", "Lectura", "Cine"],
  },
  {
    id: "4",
    name: "Isabella",
    age: 25,
    bio: "Artista y soñadora. Creo que la vida es mejor con colores y buena compañía.",
    images: [
      "https://images.unsplash.com/photo-1438761681033-6461ffad8d80?w=800&q=80",
      "https://images.unsplash.com/photo-1544005313-94ddf0286df2?w=800&q=80",
    ],
    distance: 2,
    occupation: "Artista Visual",
    interests: ["Arte", "Música", "Vino", "Teatro"],
  },
  {
    id: "5",
    name: "Mariana",
    age: 27,
    bio: "Abogada de profesión, aventurera de corazón. Siempre lista para una nueva experiencia.",
    images: [
      "https://images.unsplash.com/photo-1502823403499-6ccfcf4fb453?w=800&q=80",
      "https://images.unsplash.com/photo-1529626455594-4ff0802cfb7e?w=800&q=80",
    ],
    distance: 6,
    occupation: "Abogada",
    interests: ["Viajes", "Deportes", "Gastronomía", "Fotografía"],
  },
]

export const mockMatches: Match[] = [
  {
    id: "m1",
    profile: {
      id: "10",
      name: "Andrea",
      age: 25,
      bio: "Amante de la música y los conciertos",
      images: ["https://images.unsplash.com/photo-1531746020798-e6953c6e8e04?w=800&q=80"],
      distance: 4,
      occupation: "Músico",
      interests: ["Música", "Conciertos"],
    },
    matchedAt: new Date(Date.now() - 1000 * 60 * 30),
    lastMessage: "¿Te gustaría ir a tomar un café?",
    unread: true,
  },
  {
    id: "m2",
    profile: {
      id: "11",
      name: "Lucía",
      age: 27,
      bio: "Foodie y viajera",
      images: ["https://images.unsplash.com/photo-1524638431109-93d95c968f03?w=800&q=80"],
      distance: 7,
      occupation: "Chef",
      interests: ["Cocina", "Viajes"],
    },
    matchedAt: new Date(Date.now() - 1000 * 60 * 60 * 2),
    lastMessage: "Hola! ¿Cómo estás?",
    unread: false,
  },
  {
    id: "m3",
    profile: {
      id: "12",
      name: "Carolina",
      age: 24,
      bio: "Diseñadora y amante del arte",
      images: ["https://images.unsplash.com/photo-1485893086445-ed75865251e0?w=800&q=80"],
      distance: 3,
      occupation: "Diseñadora",
      interests: ["Arte", "Diseño"],
    },
    matchedAt: new Date(Date.now() - 1000 * 60 * 60 * 24),
    lastMessage: undefined,
    unread: false,
  },
]

export const mockMessages: Message[] = [
  {
    id: "msg1",
    senderId: "10",
    text: "¡Hola! Vi que también te gusta la música indie",
    timestamp: new Date(Date.now() - 1000 * 60 * 60),
    read: true,
  },
  {
    id: "msg2",
    senderId: "me",
    text: "¡Sí! Es mi género favorito. ¿Cuál es tu banda favorita?",
    timestamp: new Date(Date.now() - 1000 * 60 * 45),
    read: true,
  },
  {
    id: "msg3",
    senderId: "10",
    text: "Me encanta Arctic Monkeys y The Strokes",
    timestamp: new Date(Date.now() - 1000 * 60 * 40),
    read: true,
  },
  {
    id: "msg4",
    senderId: "me",
    text: "¡Excelente gusto! También son de mis favoritas",
    timestamp: new Date(Date.now() - 1000 * 60 * 35),
    read: true,
  },
  {
    id: "msg5",
    senderId: "10",
    text: "¿Te gustaría ir a tomar un café?",
    timestamp: new Date(Date.now() - 1000 * 60 * 30),
    read: false,
  },
]
