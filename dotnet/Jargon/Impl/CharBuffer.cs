using System;
using System.Collections.Generic;
using System.Text;

namespace Jargon.Impl
{
    internal struct CharBuffer
    {
        private const int INITIAL_SIZE = 8;

        private char[] Buffer;
        private int NextIndex;

        public int Length => NextIndex;

        public void Append(Rune rune)
        {
            var neededSize = NextIndex + (rune.IsMultiCharacter ? 2 : 1);

            if(neededSize > Buffer.Length)
            {
                Array.Resize(ref Buffer, NextSize(Buffer.Length));
            }

            if (rune.IsMultiCharacter)
            {
                Buffer[NextIndex] = rune.HighSurrogate;
                Buffer[NextIndex + 1] = rune.LowSurrogate;
                NextIndex += 2;
            }
            else
            {
                Buffer[NextIndex] = rune.Character;
                NextIndex++;
            }
        }

        public void Reset()
        {
            NextIndex = 0;
        }

        public string ToStringReset()
        {
            var ret = new string(Buffer, 0, NextIndex);
            NextIndex = 0;

            return ret;
        }

        public bool IsSingleRune(out Rune singleRune)
        {
            if (Length == 0)
            {
                singleRune = default;
                return false;
            }
            if (Length == 1)
            {
                singleRune = new Rune(Buffer[0]);
                return true;
            }
            if (Length > 2)
            {
                singleRune = default;
                return false;
            }

            var ret = char.IsHighSurrogate(Buffer[0]) && char.IsLowSurrogate(Buffer[1]);
            if (ret)
            {
                singleRune = new Rune(Buffer[0], Buffer[1]);
                return true;
            }

            singleRune = default;
            return false;
        }

        public static CharBuffer New()
        {
            return 
                new CharBuffer
                {
                    Buffer = new char[INITIAL_SIZE],
                    NextIndex = 0,
                };
        }

        private static int NextSize(int oldSize)
        {
            if (oldSize >= 4096)
            {
                return oldSize + 4096;
            }

            return oldSize * 2;
        }
    }
}
