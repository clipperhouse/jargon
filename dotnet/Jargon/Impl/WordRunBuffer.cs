
namespace Jargon.Impl
{
    internal struct WordRunBuffer
    {
        private string[] Buffer;
        private int LastIndex;

        public int Count => LastIndex;

        public string[] AsArray => Buffer;

        public WordRunBuffer(int maxSize)
        {
            Buffer = new string[maxSize];
            LastIndex = 0;
        }

        public void Add(string str)
        {
            Buffer[LastIndex] = str;
            LastIndex++;
        }

        public void Clear()
        {
            LastIndex = 0;
        }
    }
}
