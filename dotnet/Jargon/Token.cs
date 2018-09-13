using System;

namespace Jargon
{
    public readonly struct Token
    {
        [Flags]
        internal enum State: byte
        {
            None = 0,
            Punct = 1 << 0,
            Space = 1 << 1,
            Lemma = 1 << 2
        }

        private readonly string _String;
        public string String => _String;

        private readonly State _State;

        public bool IsPunct => _State.HasFlag(State.Punct);
        public bool IsSpace => _State.HasFlag(State.Space);
        public bool IsLemma => _State.HasFlag(State.Lemma);

        internal Token(string value, State state)
        {
            _String = value;
            _State = state;
        }

        internal Token(string value, bool punct, bool space, bool lemma)
        {
            _String = value;

            _State =
                (punct ? State.Punct : State.None) |
                (space ? State.Space : State.None) |
                (lemma ? State.Lemma : State.None);
        }

        public static Token NewLemma(string value)
        => new Token(value, State.Lemma);

        public static Token NewNone(string value)
        => new Token(value, State.None);

        public static Token NewPunct(string value)
        => new Token(value, State.Punct);

        public static Token NewSpace(string value)
        => new Token(value, State.Space);


        public override string ToString() => String;
    }
}
