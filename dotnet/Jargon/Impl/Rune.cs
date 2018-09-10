using System;

namespace Jargon.Impl
{
    internal readonly struct Rune : IEquatable<Rune>
    {
        public bool IsMultiCharacter => CodePoint > 0x0000_FFFF;    // is the high part set, basically

        public string AsString
        {
            get
            {
                if (IsMultiCharacter)
                {
                    return char.ConvertFromUtf32(char.ConvertToUtf32(HighSurrogate, LowSurrogate));
                }

                return Character.ToString();
            }
        }

        public char Character
        {
            get
            {
                if (IsMultiCharacter) throw new InvalidOperationException($"Cannot access multi-char {nameof(Rune)} with {nameof(Character)}");

                return (char)CodePoint;
            }
        }

        public char HighSurrogate
        {
            get
            {
                if (!IsMultiCharacter) throw new InvalidOperationException($"Cannot access single-char {nameof(Rune)} with {nameof(HighSurrogate)}");

                return (char)(CodePoint >> 16);
            }
        }
        public char LowSurrogate
        {
            get
            {
                if (!IsMultiCharacter) throw new InvalidOperationException($"Cannot access single-char {nameof(Rune)} with {nameof(LowSurrogate)}");

                return (char)(CodePoint & 0xFFFF);
            }
        }

        private readonly uint CodePoint;


        public Rune(char c)
        {
            if (char.IsHighSurrogate(c) || char.IsLowSurrogate(c)) throw new ArgumentException($"Cannot create single-char {nameof(Rune)} with surrogate char", nameof(c));

            CodePoint = c;
        }

        public Rune(char high, char low)
        {
            if (!char.IsHighSurrogate(high)) throw new ArgumentException($"Cannot create multi-char {nameof(Rune)} with non-high-surrogate char", nameof(high));
            if (!char.IsLowSurrogate(low)) throw new ArgumentException($"Cannot create multi-char {nameof(Rune)} with non-high-surrogate char", nameof(low));

            CodePoint = (uint)((high << 16) | (low << 0));
        }

        public bool IsSpace()
        {
            if (IsMultiCharacter)
            {
                return char.IsWhiteSpace(AsString, 0);
            }

            return char.IsWhiteSpace(Character);
        }

        public bool MightBeLeadingPunct()
        {
            if (IsMultiCharacter) return false;

            // only case
            return Character == '.';
        }

        public bool MightBeMidPunct()
        {
            if (IsMultiCharacter) return false;

            switch (Character)
            {
                case '.':
                case '\'':
                case '’':
                case ':':
                case '?':
                case '&':
                    return true;
            }

            return false;
        }

        public bool IsPunct()
        {
            if (IsMultiCharacter)
            {
                return char.IsPunctuation(AsString, 0); // none of the exceptions are multi-char, so elide the call
            }

            var c = Character;
            return char.IsPunctuation(c) && !IsPunctException(c);

            bool IsPunctException(char cp)
            {
                switch (cp)
                {
                    case '-':
                    case '#':
                    case '@':
                    case '*':
                    case '%':
                    case '_':
                    case '/':
                    case '\\': return true;
                }

                return false;
            }
        }
        
        public bool Equals(Rune other) => other.CodePoint == CodePoint;

        public override bool Equals(object obj)
        {
            if (obj is Rune r)
            {
                return Equals(r);
            }

            return false;
        }

        public override int GetHashCode() => (int)CodePoint;

        public override string ToString() => AsString;
    }
}
