-- migrate:up
create or replace function uuid7() returns uuid as $$
declare
begin
	return uuid7(clock_timestamp());
end $$ language plpgsql;

create or replace function uuid7(p_timestamp timestamp with time zone) returns uuid as $$
declare

	v_time double precision := null;

	v_unix_t bigint := null;
	v_rand_a bigint := null;
	v_rand_b bigint := null;

	v_unix_t_hex varchar := null;
	v_rand_a_hex varchar := null;
	v_rand_b_hex varchar := null;

	c_milli double precision := 10^3;  -- 1 000
	c_micro double precision := 10^6;  -- 1 000 000
	c_scale double precision := 4.096; -- 4.0 * (1024 / 1000)

	c_version bigint := x'0000000000007000'::bigint; -- RFC-9562 version: b'0111...'
	c_variant bigint := x'8000000000000000'::bigint; -- RFC-9562 variant: b'10xx...'

begin

	v_time := extract(epoch from p_timestamp);

	v_unix_t := trunc(v_time * c_milli);
	v_rand_a := trunc((v_time * c_micro - v_unix_t * c_milli) * c_scale);
	v_rand_b := trunc(random() * 2^30)::bigint << 32 | trunc(random() * 2^32)::bigint;

	v_unix_t_hex := lpad(to_hex(v_unix_t), 12, '0');
	v_rand_a_hex := lpad(to_hex((v_rand_a | c_version)::bigint), 4, '0');
	v_rand_b_hex := lpad(to_hex((v_rand_b | c_variant)::bigint), 16, '0');

	return (v_unix_t_hex || v_rand_a_hex || v_rand_b_hex)::uuid;

end $$ language plpgsql;

-- Base62 encoding helper function
create or replace function base62_encode(p_number bigint) returns text as $$
declare
    v_chars text := '0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz';
    v_result text := '';
    v_num bigint := abs(p_number);
begin
    if v_num = 0 then
        return '0';
    end if;

    while v_num > 0 loop
        v_result := substr(v_chars, (v_num % 62)::integer + 1, 1) || v_result;
        v_num := v_num / 62;
    end loop;

    return v_result;
end;
$$ language plpgsql;

create or replace function generate_api_key(p_prefix text) returns text as $$
declare
    v_result text;
    v_uuid uuid;
    v_high bigint;
    v_low bigint;
begin
    -- Validate prefix
    if p_prefix not in ('sk_live_', 'pk_live_', 'sk_test_', 'pk_test_') then
        raise exception 'Invalid API key prefix. Must be one of: sk_live_, pk_live_, sk_test_, pk_test_';
    end if;

    -- Generate UUID v7 for timestamp-ordered uniqueness
    v_uuid := uuid7();

    -- Split UUID into two 64-bit integers
    v_high := (('x' || substr(v_uuid::text, 1, 16))::bit(64))::bigint;
    v_low := (('x' || substr(v_uuid::text, 17, 16))::bit(64))::bigint;

    -- Combine prefix with base62 encoded values
    v_result := p_prefix || base62_encode(v_high) || base62_encode(v_low);

    return v_result;
end;
$$ language plpgsql;

-- migrate:down
drop function if exists uuid7(timestamp with time zone);
drop function if exists uuid7();
drop function if exists generate_api_key(text);
drop function if exists base62_encode(bigint);
