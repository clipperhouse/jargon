using System;

namespace Jargon.Impl
{
    internal struct TokenBuffer
    {
        private Token[] Buffer;
        private int FirstIndex;
        private int NextIndex;

        public int Count
        {
            get
            {
                if (FirstIndex == -1) return 0;

                return NextIndex - FirstIndex;
            }
        }

        public TokenBuffer(int initialSize)
        {
            Buffer = new Token[initialSize];
            FirstIndex = -1;
            NextIndex = 0;
        }

        public Token FirstElement() => Buffer[FirstIndex];

        public Token ElementAt(int ix) => Buffer[FirstIndex + ix];

        public void Add(Token token)
        {
            // 3 cases
            //   1) NextIndex != Buffer.Length
            //      => if FirstIndex == -1, FirstIndex = NextIndex
            //      => put token in Buffer, incr NextIndex
            //   2) NextIndex == Buffer.Length, FirstIndex == 0
            //      => increase the size of Buffer
            //      => put token in Buffer, incr NextIndex
            //   3) NextIndex == Buffer.Length, FirstIndex > 0
            //      => copy buffer back into itself
            //      => NextIndex -= FirstIndex
            //      => FirstIndex = 0
            //      => put token in Buffer, incr NextIndex

            if(NextIndex != Buffer.Length)
            {
                if(FirstIndex == -1)
                {
                    FirstIndex = NextIndex;
                }
                Buffer[NextIndex] = token;
                NextIndex++;
                return;
            }

            if(NextIndex == Buffer.Length && FirstIndex == 0)
            {
                Array.Resize(ref Buffer, NextSize(Buffer.Length));
                Buffer[NextIndex] = token;
                NextIndex++;
                return;
            }

            Array.Copy(Buffer, FirstIndex, Buffer, 0, Count);
            NextIndex -= FirstIndex;
            FirstIndex = 0;
            Buffer[NextIndex] = token;
            NextIndex++;
        }

        public void RemoveFromFront(int count)
        {
            FirstIndex += count;
            FirstIndex = Math.Min(NextIndex, FirstIndex);
        }

        private static int NextSize(int oldSize)
        {
            if(oldSize >= 4098)
            {
                return oldSize + 4098; 
            }

            return oldSize * 2;
        }
    }
}
