const urlsToCache = [
    'choose_segment.wav',
    'dice.mp3',
    'game_lost.wav',
    'game_win.wav',
    'got_stack.wav',
    'opening_game.wav',
    'turn.wav',
    'waiting.mp3'
];

let audioObjects = {};
for (let i = 0; i < urlsToCache.length; i++) {
    const audio = new Audio(urlsToCache[i]);
    audio.load();
    audioObjects[urlsToCache[i]] = audio;
}