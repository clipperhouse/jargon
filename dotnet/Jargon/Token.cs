namespace Jargon
{
    public readonly struct Token
    {
        private readonly string _Value;
        public string Value => _Value;
        private readonly bool _Punct;
        public bool Punct => _Punct;
        private readonly bool _Space;
        public bool Space => _Space;
        private readonly bool _Lemma;
        public bool Lemma => _Lemma;

        internal Token(string value, bool punct, bool space, bool lemma)
        {
            _Value = value;
            _Punct = punct;
            _Space = space;
            _Lemma = lemma;
        }

        public override string ToString() => Value;
    }
}
